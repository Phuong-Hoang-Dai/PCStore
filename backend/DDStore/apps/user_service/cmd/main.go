package main

import (
	"log"

	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/db"
	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/handler"
	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/logger"
)

func main() {
	logger.Setup("user_service")
	configs.LoadConfig()

	db, err := db.SetupPostgres()

	if err != nil {
		log.Fatal("Error connecting database")
	}

	handler.SetupHttp(db)
}
