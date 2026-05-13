package db

import (
	"context"
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/configs"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func SetupDB() (client *mongo.Client, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	dsn := configs.Cfg.ConnectStr
	client, err = mongo.Connect(options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
