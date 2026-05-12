package main

import (
	"log"

	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/db"
	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/handler"
)

func main() {
	configs.LoadConfig()

	db, err := db.SetupDB()

	if err != nil {
		log.Fatal("Error connecting database")
	}

	handler.SetupHttp(db)
}
