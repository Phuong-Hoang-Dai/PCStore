package repos

import (
	"context"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgPaymentRepo struct {
	DB *pgxpool.Pool
}

func NewPgPaymentRepo(db *pgxpool.Pool) service.PaymentRepository {
	return pgPaymentRepo{
		DB: db,
	}
}

const orderPaymentColumns = `order_id, user_id, state, amount_total, currency,
	stripe_user_id, payment_intent_id, checkout_session_id`

func (p pgPaymentRepo) CreateOrderPayment(data model.OrderPayment, ctx context.Context) (int, error) {
	if data.PaymentInfo == nil {
		data.PaymentInfo = &model.Payment{}
	}
	var id int
	err := p.DB.QueryRow(ctx,
		`INSERT INTO order_payments (`+orderPaymentColumns+`)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`,
		data.OrderId, data.UserId, data.State, data.AmountTotal, data.Currency,
		data.PaymentInfo.StripeUserId, data.PaymentInfo.PaymentIntentId, data.PaymentInfo.CheckoutSessionId,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p pgPaymentRepo) UpdateOrderPayment(data model.OrderPayment, ctx context.Context) error {
	result, err := p.DB.Exec(ctx,
		`UPDATE order_payments SET state = $2, stripe_user_id = $3, payment_intent_id = $4,
		checkout_session_id = $5, updated_at = NOW()
		WHERE order_id = $1`,
		data.OrderId, data.State,
		data.PaymentInfo.StripeUserId, data.PaymentInfo.PaymentIntentId, data.PaymentInfo.CheckoutSessionId)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// Stripe delivers webhooks at least once and in no guaranteed order, so both methods below
// are safe to run for a duplicate or late event:
//   - the event ID is recorded in the same transaction, and an already-recorded event is skipped
//   - state only moves forward (a completed payment / paid order is never overwritten)

func (p pgPaymentRepo) CompletePayment(eventId string, orderId int, sessionId string, ctx context.Context) error {
	return p.runOnce(eventId, ctx, func(tx pgx.Tx) error {
		// Any paid session counts, even if it isn't the latest one stored: the customer was charged.
		_, err := tx.Exec(ctx,
			`UPDATE order_payments SET state = $2, checkout_session_id = $3, updated_at = NOW()
			WHERE order_id = $1 AND state <> $2`,
			orderId, model.IsCompleted, sessionId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx,
			`UPDATE orders SET state = $2, updated_at = NOW()
			WHERE id = $1 AND state IN ($3, $4) AND deleted_at IS NULL`,
			orderId, model.Paid, model.Processing, model.ReadyToPay)
		return err
	})
}

func (p pgPaymentRepo) FailPayment(eventId string, orderId int, sessionId string, ctx context.Context) error {
	return p.runOnce(eventId, ctx, func(tx pgx.Tx) error {
		// Only the session currently stored for the order can fail it, so an older abandoned
		// session expiring can't override a newer one.
		_, err := tx.Exec(ctx,
			`UPDATE order_payments SET state = $2, updated_at = NOW()
			WHERE order_id = $1 AND checkout_session_id = $3 AND state <> $4`,
			orderId, model.IsFailed, sessionId, model.IsCompleted)
		return err
	})
}

// runOnce records the Stripe event ID and runs fn in one transaction, so an event is marked as
// handled only if its changes commit. An event that was already recorded is skipped.
func (p pgPaymentRepo) runOnce(eventId string, ctx context.Context, fn func(tx pgx.Tx) error) error {
	return pgx.BeginFunc(ctx, p.DB, func(tx pgx.Tx) error {
		result, err := tx.Exec(ctx,
			`INSERT INTO stripe_webhook_events (event_id) VALUES ($1) ON CONFLICT DO NOTHING`, eventId)
		if err != nil {
			return err
		}
		if result.RowsAffected() == 0 {
			return nil
		}
		return fn(tx)
	})
}

func (p pgPaymentRepo) GetOrderPayment(orderId int, ctx context.Context) (data model.OrderPayment, err error) {
	data.PaymentInfo = &model.Payment{}
	err = p.DB.QueryRow(ctx,
		`SELECT `+orderPaymentColumns+` FROM order_payments WHERE order_id = $1`, orderId,
	).Scan(&data.OrderId, &data.UserId, &data.State, &data.AmountTotal, &data.Currency,
		&data.PaymentInfo.StripeUserId, &data.PaymentInfo.PaymentIntentId, &data.PaymentInfo.CheckoutSessionId)
	return data, err
}
