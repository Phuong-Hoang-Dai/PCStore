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
)

func main() {
	configs.LoadConfig()

	esClient, err := db.SetupES()
	if err != nil {
		log.Fatal("Error connecting elasticsearch: ", err)
	}

	srv := handler.SetupHttp(esClient)

	go func() {
		log.Printf("search_service listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Error starting server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("forced shutdown:", err)
	}
	log.Println("done")
}
