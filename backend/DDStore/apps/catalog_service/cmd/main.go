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
	hl "github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/handler"
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

	handler := hl.NewRabbitHandler(mongoClient, redisClient, ctx)
	go handler.ProductCreatedConsumer(configs.Cfg.ProductCreatedQueue, handler.ProductCreatedWorker, context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	redisClient.Close()
	mongoClient.Disconnect(shutdownCtx)
	log.Println("done")
}
