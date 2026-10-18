package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupHttp(db *pgxpool.Pool) {
	baseUrl := "http://product_service:8888"

	orderService := Init(db, baseUrl)

	r := gin.Default()

	order := r.Group("/order")
	{
		order.POST("/:id", orderService.CreateOrder())
		order.GET("/:id", orderService.GetOrderById())
		order.GET("", orderService.GetOrders())
		order.GET("/history/:id", orderService.GetHistoryOrders())
		order.DELETE("/:id", orderService.CancelOrder())
		order.PUT("/:id", orderService.UpdateOrderState())
		order.POST("/completeOrder/:id", orderService.CompleteOrder())
		order.GET("/healthcheck", func(ctx *gin.Context) {
			ctx.JSON(http.StatusAccepted, gin.H{
				"message": "order service is healthy",
			})
		})
	}

	r.Run(":7777")
}
