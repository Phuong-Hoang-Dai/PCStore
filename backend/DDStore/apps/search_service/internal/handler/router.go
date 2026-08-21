package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/configs"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
)

// SetupHttp wires the dependencies, ensures the index exists and returns a
// ready-to-run *http.Server (started/stopped by the caller for graceful shutdown).
func SetupHttp(es *elasticsearch.TypedClient) *http.Server {
	searchHandler := NewSearchHandler(es)

	if _, err := searchHandler.service.EnsureIndex(context.Background()); err != nil {
		// Non-fatal: indexing will lazily create the index with a dynamic
		// mapping if needed. Log so the operator can investigate.
		log.Printf("warning: ensure index failed: %v", err)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	products := r.Group("/search/products")
	{
		products.POST("", searchHandler.IndexProduct())       // insert one product
		products.POST("/bulk", searchHandler.IndexProducts()) // insert many products
		products.DELETE("/:id", searchHandler.DeleteProduct())
		products.POST("/search", searchHandler.Search())             // full-text search + filters
		products.POST("/autocomplete", searchHandler.AutoComplete()) // autocomplete
	}

	return &http.Server{
		Addr:    ":" + configs.Cfg.Port,
		Handler: r,
	}
}
