package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/db"
	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/internal/handler"
	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/internal/repos"
	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/internal/service"
)

func main() {
	configs.LoadConfig()

	esClient, err := db.SetupES()
	if err != nil {
		log.Fatal("Error connecting elasticsearch: ", err)
	}

	searchService := service.NewSearchService(repos.NewESProductRepo(esClient, configs.Cfg.ESIndex))
	srv := handler.SetupHttp(searchService)

	go func() {
		log.Printf("search_service listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Error starting server: ", err)
		}
	}()

	consumerCtx, stopConsumer := context.WithCancel(context.Background())
	go func() {
		if err := handler.NewProductIndexConsumer(searchService).Run(consumerCtx); err != nil {
			log.Fatal("Error consuming product index queue: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	stopConsumer()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("forced shutdown:", err)
	}
	log.Println("done")
}
