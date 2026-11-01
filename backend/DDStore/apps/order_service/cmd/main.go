package main

import (
	"context"
	"log"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/db"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/handler"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/logger"
)

func main() {
	logger.Setup("order_service")
	configs.LoadConfig()

	// Before SetupPostgres: its result is named db and hides the package.
	rdb, err := db.SetupRedis(context.Background())
	if err != nil {
		log.Fatal("Error connecting redis: ", err)
	}

	db, err := db.SetupPostgres()

	if err != nil {
		log.Fatal("Error connecting database")
	}

	handler.SetupHttp(db, rdb)
}
