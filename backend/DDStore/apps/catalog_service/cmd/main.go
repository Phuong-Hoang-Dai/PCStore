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
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/repos"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/service"
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

	conn := hl.SetupRabbitQueues(ctx)
	publisher := hl.NewRabbitPublisher(conn, ctx)

	database := mongoClient.Database(configs.Cfg.DBName)
	productRepo := repos.NewMongoProductRepo(database.Collection("products"))
	cateRepo := repos.NewMongoCateRepo(database)
	productService := service.NewProductService(productRepo, cateRepo, redisClient, publisher)
	cateService := service.NewCateService(cateRepo, productRepo, redisClient, publisher)

	seeded, err := db.SeedData(ctx, database)
	if err != nil {
		log.Fatal("Error seeding database: ", err)
	}
	// Seeded products bypass ProductService, so index them for search here.
	for _, p := range seeded {
		if err := publisher.ProductUpserted(ctx, p); err != nil {
			log.Printf("Failed to publish seeded product %s: %v", p.Id.Hex(), err)
		}
	}

	consumer := hl.NewRabbitConsumer(conn, productService)
	go consumer.CreateConsumer(configs.Cfg.ProductCreatedQueue, consumer.ProductCreatedWorker, context.Background())
	go hl.SetupHttp(productService, cateService)

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
