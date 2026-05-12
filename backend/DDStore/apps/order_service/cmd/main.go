package main

import (
	"log"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/db"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/handler"
)

func main() {
	configs.LoadConfig()

	db, err := db.SetupDB()

	if err != nil {
		log.Fatal("Error connecting database")
	}

	handler.SetupHttp(db)
}
