package db

import (
	"context"
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/configs"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func SetupDB(ctx context.Context) (*mongo.Client, *mongo.Database, error) {
	opts := options.Client().
		ApplyURI(configs.Cfg.MongoURI).
		SetMaxPoolSize(10).
		SetMinPoolSize(2).
		SetMaxConnIdleTime(5 * time.Minute)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, nil, err
	}

	return client, client.Database(configs.Cfg.DBName), nil
}