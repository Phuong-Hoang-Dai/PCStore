package db

import (
	"context"
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/configs"
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
	if err := SeedData(pool); err != nil {
		return pool, err
	}
	return pool, nil
}

func Migrate(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id         SERIAL PRIMARY KEY,
			name       VARCHAR(255) NOT NULL,
			email      VARCHAR(255) NOT NULL DEFAULT '',
			password   BYTEA,
			role       VARCHAR(50) NOT NULL DEFAULT 'user',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ
		);
		CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);
		CREATE INDEX IF NOT EXISTS idx_users_name ON users (name);
		-- Guest accounts (created by orders placed without logging in) have no password.
		ALTER TABLE users ALTER COLUMN password DROP NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_users_email ON users (lower(email));
	`)
	return err
}
