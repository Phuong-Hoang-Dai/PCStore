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
	return pool, nil
}
