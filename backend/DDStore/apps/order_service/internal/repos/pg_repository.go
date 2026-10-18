package repos

import (
	"context"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgOrderRepo struct {
	DB *pgxpool.Pool
}

func NewPgOrderRepo(db *pgxpool.Pool) service.OrderRepository {
	return pgOrderRepo{
		DB: db,
	}
}

const orderColumns = `id, state, total, user_id, created_at, updated_at`

func (m pgOrderRepo) CreateOrder(data *model.Order, ctx context.Context) (int, error) {
	err := pgx.BeginFunc(ctx, m.DB, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			`INSERT INTO orders (state, total, user_id) VALUES ($1, $2, $3)
			RETURNING id, created_at, updated_at`,
			data.State, data.Total, data.UserId,
		).Scan(&data.Id, &data.CreatedAt, &data.UpdatedAt)
		if err != nil {
			return err
		}

		_, err = tx.CopyFrom(ctx,
			pgx.Identifier{"order_items"},
			[]string{"order_id", "product_id", "quantity", "price"},
			pgx.CopyFromSlice(len(data.Items), func(i int) ([]any, error) {
				data.Items[i].OrderId = data.Id
				item := data.Items[i]
				return []any{item.OrderId, item.ProductId, item.Quantity, item.Price}, nil
			}),
		)
		return err
	})
	if err != nil {
		return 0, err
	}

	return data.Id, nil
}

func (m pgOrderRepo) GetOrderById(id int, ctx context.Context) (data model.Order, err error) {
	err = m.DB.QueryRow(ctx,
		`SELECT `+orderColumns+` FROM orders WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&data.Id, &data.State, &data.Total, &data.UserId, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		return data, err
	}

	rows, err := m.DB.Query(ctx,
		`SELECT product_id, quantity, order_id, price, created_at, updated_at
		FROM order_items WHERE order_id = $1`, id)
	if err != nil {
		return data, err
	}

	data.Items, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (item model.OrderItem, err error) {
		err = row.Scan(&item.ProductId, &item.Quantity, &item.OrderId, &item.Price, &item.CreatedAt, &item.UpdatedAt)
		return item, err
	})

	return data, err
}

func (m pgOrderRepo) UpdateOrder(data model.Order, ctx context.Context) error {
	result, err := m.DB.Exec(ctx,
		`UPDATE orders SET state = $2, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`,
		data.Id, data.State)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (m pgOrderRepo) DeleteOrder(id int, ctx context.Context) error {
	result, err := m.DB.Exec(ctx,
		`UPDATE orders SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// LIMIT NULLIF($1, 0): a limit of 0 means no limit (LIMIT NULL).
func (m pgOrderRepo) GetOrders(p model.Paging, ctx context.Context) (orders []model.Order, err error) {
	rows, err := m.DB.Query(ctx,
		`SELECT `+orderColumns+` FROM orders WHERE deleted_at IS NULL
		ORDER BY id LIMIT NULLIF($1, 0) OFFSET $2`,
		p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}

	return collectOrders(rows)
}

func (m pgOrderRepo) GetHistoryOrders(id int, p model.Paging, ctx context.Context) (orders []model.Order, err error) {
	rows, err := m.DB.Query(ctx,
		`SELECT `+orderColumns+` FROM orders WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT NULLIF($2, 0) OFFSET $3`,
		id, p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}

	return collectOrders(rows)
}

func collectOrders(rows pgx.Rows) ([]model.Order, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (o model.Order, err error) {
		err = row.Scan(&o.Id, &o.State, &o.Total, &o.UserId, &o.CreatedAt, &o.UpdatedAt)
		return o, err
	})
}
