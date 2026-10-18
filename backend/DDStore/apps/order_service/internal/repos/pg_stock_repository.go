package repos

import (
	"context"
	"errors"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgStockRepo struct {
	DB *pgxpool.Pool
}

func NewPgStockRepo(db *pgxpool.Pool) service.StockRepository {
	return pgStockRepo{
		DB: db,
	}
}

func (p pgStockRepo) ReserveStock(productId string, quantity int, orderId int, ctx context.Context) error {
	err := pgx.BeginFunc(ctx, p.DB, func(tx pgx.Tx) (err error) {
		stockId := 0
		err = tx.QueryRow(ctx,
			`UPDATE stock SET quantity = quantity - $2, updated_at = NOW() 
			WHERE product_id = $1 AND quantity >= $2 AND deleted_at IS NULL
			RETURNING stock_id`, productId, quantity,
		).Scan(&stockId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO stock_reservation (product_id, order_id, quantity, stock_id) 
			VALUES ($1, $2, $3, $4)`, productId, orderId, quantity, stockId,
		)
		if err != nil {
			return err
		}
		return nil
	})

	if errors.Is(err, pgx.ErrNoRows) {
		err = model.ErrOutOfStock
	}

	return err
}
