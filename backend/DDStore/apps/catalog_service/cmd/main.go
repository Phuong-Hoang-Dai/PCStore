package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/db"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/handler"
)

func main() {
	configs.LoadConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, _, err := db.SetupDB(ctx)
	if err != nil {
		log.Fatal("Error connecting database: ", err)
	}

	redisClient, err := db.SetupRedis(ctx)
	if err != nil {
		log.Fatal("Error connecting redis: ", err)
	}

	handler.SetupHttp(mongoClient, redisClient)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	mongoClient.Disconnect(shutdownCtx)
	redisClient.Close()
	log.Println("done")
}
