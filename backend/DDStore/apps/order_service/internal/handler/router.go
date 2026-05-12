package handler

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func SetupHttp(db *gorm.DB) {
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
	}

	r.Run(":7777")
}
