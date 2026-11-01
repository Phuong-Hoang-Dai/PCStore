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

// ReserveStock takes quantity out of stock for the order and returns what is left.
// ErrOutOfStock means the reservation was made and took the last units.
// ErrNotEnoughStock comes with the quantity in stock, which was left untouched.
func (p pgStockRepo) ReserveStock(productId string, quantity int, orderId int, ctx context.Context) (remaining int, err error) {
	err = pgx.BeginFunc(ctx, p.DB, func(tx pgx.Tx) (err error) {
		stockId := 0
		err = tx.QueryRow(ctx,
			`UPDATE stock SET quantity = quantity - $2, updated_at = NOW()
			WHERE product_id = $1 AND quantity >= $2 AND deleted_at IS NULL
			RETURNING stock_id, quantity`, productId, quantity,
		).Scan(&stockId, &remaining)
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
		// A product without a stock row has none.
		err = p.DB.QueryRow(ctx,
			`SELECT quantity FROM stock WHERE product_id = $1 AND deleted_at IS NULL`, productId,
		).Scan(&remaining)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return 0, err
		}
		return remaining, model.ErrNotEnoughStock
	}
	if err != nil {
		return 0, err
	}

	if remaining == 0 {
		return 0, model.ErrOutOfStock
	}
	return remaining, nil
}

// ReleaseStock returns every active reservation of the order to stock and returns the
// restocked product ids.
func (p pgStockRepo) ReleaseStock(orderId int, ctx context.Context) ([]string, error) {
	return releaseReservations(ctx, p.DB, []int{orderId})
}

// AddStock adds quantity to a product's stock, creating its stock row if needed,
// and returns the new quantity.
func (p pgStockRepo) AddStock(productId string, quantity int, ctx context.Context) (total int, err error) {
	err = p.DB.QueryRow(ctx,
		`INSERT INTO stock (product_id, quantity) VALUES ($1, $2)
		ON CONFLICT (product_id) DO UPDATE
		SET quantity = stock.quantity + EXCLUDED.quantity, updated_at = NOW(), deleted_at = NULL
		RETURNING quantity`, productId, quantity,
	).Scan(&total)
	return total, err
}

// GetStocks returns the stock of the given products, or of all products when ids is empty.
func (p pgStockRepo) GetStocks(productIds []string, ctx context.Context) ([]model.Stock, error) {
	var filter []string
	if len(productIds) > 0 {
		filter = productIds
	}
	rows, err := p.DB.Query(ctx,
		`SELECT product_id, quantity, updated_at FROM stock
		WHERE deleted_at IS NULL AND ($1::text[] IS NULL OR product_id = ANY($1))
		ORDER BY product_id`, filter)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (s model.Stock, err error) {
		err = row.Scan(&s.ProductId, &s.Quantity, &s.UpdatedAt)
		return s, err
	})
}

type querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// releaseReservations gives the active reservations of the orders back to stock, in one
// statement so a reservation can't be released twice, and returns the restocked product ids.
// It takes a querier so order cancellation can run it in the same transaction.
func releaseReservations(ctx context.Context, q querier, orderIds []int) ([]string, error) {
	rows, err := q.Query(ctx,
		`WITH released AS (
			UPDATE stock_reservation SET deleted_at = NOW(), updated_at = NOW()
			WHERE order_id = ANY($1) AND deleted_at IS NULL
			RETURNING stock_id, quantity
		)
		UPDATE stock s SET quantity = s.quantity + r.quantity, updated_at = NOW()
		FROM (SELECT stock_id, SUM(quantity) AS quantity FROM released GROUP BY stock_id) r
		WHERE s.stock_id = r.stock_id
		RETURNING s.product_id`, orderIds)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}
