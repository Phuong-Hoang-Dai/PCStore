package main

import (
	"fmt"
	"log"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/db"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/handler"
)

func main() {
	configs.LoadConfig()
	fmt.Println("this is connectr: ", configs.Cfg.ConnectStr)

	db, err := db.SetupDB()

	if err != nil {
		log.Fatal("Error connecting database")
	}

	handler.SetupHttp(db)
}
