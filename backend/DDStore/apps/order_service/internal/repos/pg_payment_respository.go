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

func (p pgPaymentRepo) CreateOrderPayment(data model.OrderPayment, ctx context.Context) (int, error) {
	var id int
	err := p.DB.QueryRow(ctx,
		`INSERT INTO order_payments (order_id, user_id, checkout_session_id, state, payment_intent_id, amount_total, currency)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		data.OrderId, data.UserId, data.CheckoutSessionId, data.State, data.PaymentIntentId, data.AmountTotal, data.Currency,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p pgPaymentRepo) UpdateOrderPayment(data model.OrderPayment, ctx context.Context) error {
	result, err := p.DB.Exec(ctx,
		`UPDATE order_payments SET checkout_session_id = $2, state = $3, payment_intent_id = $4, updated_at = NOW()
		WHERE order_id = $1`,
		data.OrderId, data.CheckoutSessionId, data.State, data.PaymentIntentId)
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
		`SELECT order_id, user_id, checkout_session_id, state, payment_intent_id, amount_total, currency
		FROM order_payments WHERE order_id = $1`, orderId,
	).Scan(&data.OrderId, &data.UserId, &data.CheckoutSessionId, &data.State,
		&data.PaymentIntentId, &data.AmountTotal, &data.Currency)
	return data, err
}
