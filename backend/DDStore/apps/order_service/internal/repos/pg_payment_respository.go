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

func (p pgPaymentRepo) GetOrderPayment(orderId int, ctx context.Context) (data model.OrderPayment, err error) {
	err = p.DB.QueryRow(ctx,
		`SELECT `+orderPaymentColumns+` FROM order_payments WHERE order_id = $1`, orderId,
	).Scan(&data.OrderId, &data.UserId, &data.State, &data.AmountTotal, &data.Currency,
		&data.PaymentInfo.StripeUserId, &data.PaymentInfo.PaymentIntentId, &data.PaymentInfo.CheckoutSessionId)
	return data, err
}
