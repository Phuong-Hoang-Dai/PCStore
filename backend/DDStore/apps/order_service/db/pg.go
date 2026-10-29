package db

import (
	"context"
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupPostgres() (db *pgxpool.Pool, err error) {
	cfg, err := pgxpool.ParseConfig(configs.Cfg.ConnectStr)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}
	if err := Migrate(pool); err != nil {
		pool.Close()
		return nil, err
	}
	if err := SeedStock(pool); err != nil {
		return pool, err
	}
	return pool, nil
}

// orders.id has no default: ids come from order_sequencer (see pg_sequencer_repository.go).
func Migrate(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS orders (
			id         INTEGER PRIMARY KEY,
			state      INTEGER NOT NULL DEFAULT 0,
			total      NUMERIC(15, 2) NOT NULL DEFAULT 0,
			user_id    INTEGER NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ
		);
		CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders (user_id);

		CREATE TABLE IF NOT EXISTS order_items (
			order_id   INTEGER NOT NULL REFERENCES orders (id),
			product_id VARCHAR(24) NOT NULL,
			quantity   INTEGER NOT NULL,
			price      NUMERIC(15, 2) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (order_id, product_id)
		);

		CREATE TABLE IF NOT EXISTS order_sequencer (
			id            INTEGER PRIMARY KEY,
			current_value INTEGER NOT NULL,
			max_value     INTEGER NOT NULL
		);
		CREATE TABLE IF NOT EXISTS range_handler (
			id              INTEGER PRIMARY KEY,
			available_range INTEGER NOT NULL,
			usable_quantity INTEGER NOT NULL
		);
		INSERT INTO order_sequencer (id, current_value, max_value) VALUES (1, 0, 0) ON CONFLICT DO NOTHING;
		INSERT INTO range_handler (id, available_range, usable_quantity) VALUES (1, 0, 1000) ON CONFLICT DO NOTHING;

		CREATE TABLE IF NOT EXISTS stock (
			stock_id   SERIAL PRIMARY KEY,
			product_id VARCHAR(24) NOT NULL UNIQUE,
			quantity   INTEGER NOT NULL CHECK (quantity >= 0),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ
		);
		CREATE TABLE IF NOT EXISTS stock_reservation (
			id         SERIAL PRIMARY KEY,
			product_id VARCHAR(24) NOT NULL,
			order_id   INTEGER NOT NULL,
			quantity   INTEGER NOT NULL,
			stock_id   INTEGER NOT NULL REFERENCES stock (stock_id),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ
		);
		CREATE INDEX IF NOT EXISTS idx_stock_reservation_order_id ON stock_reservation (order_id);

		CREATE TABLE IF NOT EXISTS order_payments (
			id                  SERIAL PRIMARY KEY,
			order_id            INTEGER NOT NULL UNIQUE REFERENCES orders (id),
			user_id             INTEGER NOT NULL,
			state               INTEGER NOT NULL DEFAULT 0,
			amount_total        INTEGER,
			currency            VARCHAR(3),
			stripe_user_id      VARCHAR(255) NOT NULL DEFAULT '',
			payment_intent_id   VARCHAR(255) NOT NULL DEFAULT '',
			checkout_session_id VARCHAR(255) NOT NULL DEFAULT '',
			created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS stripe_webhook_events (
			event_id    VARCHAR(255) PRIMARY KEY,
			received_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	return err
}

// seedProductIds must match catalog_service/db/seed.go SeedProductIds.
// There is no stock management API yet, so only these products can be ordered.
var seedProductIds = []string{
	"665000000000000000000001",
	"665000000000000000000002",
	"665000000000000000000003",
	"665000000000000000000004",
	"665000000000000000000005",
}

const seedStockQuantity = 100

func SeedStock(db *pgxpool.Pool) error {
	for _, id := range seedProductIds {
		_, err := db.Exec(context.Background(),
			`INSERT INTO stock (product_id, quantity) VALUES ($1, $2) ON CONFLICT (product_id) DO NOTHING`,
			id, seedStockQuantity)
		if err != nil {
			return err
		}
	}
	return nil
}
