package repos

import (
	"context"
	"errors"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgSequencerkRepo struct {
	DB *pgxpool.Pool
}

func NewPgSequencerkRepo(db *pgxpool.Pool) service.SequenecerRepository {
	return pgSequencerkRepo{
		DB: db,
	}
}

func (p pgSequencerkRepo) GetOrderId(ctx context.Context) (orderId int, err error) {
	err = pgx.BeginFunc(ctx, p.DB, func(tx pgx.Tx) (err error) {
		err = tx.QueryRow(ctx,
			`UPDATE order_sequencer SET current_value = current_value + 1
			WHERE id = 1 AND current_value < max_value
			RETURNING current_value`,
		).Scan(&orderId)

		return err
	})

	if err != nil && err != pgx.ErrNoRows {
		return -1, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		orderId, err = p.getNewRange(ctx)

		if err != nil {
			return -1, err
		}
	}

	return orderId, nil
}

func (p pgSequencerkRepo) getNewRange(ctx context.Context) (orderId int, err error) {
	err = pgx.BeginFunc(ctx, p.DB, func(tx pgx.Tx) (err error) {
		newRange, usableQuantity := 0, 0
		err = tx.QueryRow(ctx,
			`UPDATE range_handler SET available_range = available_range + usable_Quantity + 1
				WHERE id = 1  
				RETURNING available_range, usable_Quantity`,
		).Scan(&newRange, &usableQuantity)

		if err != nil {
			return err
		}
		err = tx.QueryRow(ctx,
			`UPDATE order_sequencer SET current_value = $1, max_value = $1 + $2
				WHERE id = 1  
				RETURNING current_value`, newRange, usableQuantity,
		).Scan(&orderId)

		return err
	})

	if err != nil {
		return -1, err
	}

	return orderId, nil
}
