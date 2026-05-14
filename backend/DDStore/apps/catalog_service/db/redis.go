package db

import (
	"context"
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/configs"
	"github.com/redis/go-redis/v9"
)

func SetupRedis(ctx context.Context) (*redis.Client, error) {
	opts, err := redis.ParseURL(configs.Cfg.RedisAddr)
	if err != nil {
		return nil, err
	}

	opts.PoolSize = 10
	opts.MinIdleConns = 2
	opts.ConnMaxIdleTime = 5 * time.Minute

	rdb := redis.NewClient(opts)

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
