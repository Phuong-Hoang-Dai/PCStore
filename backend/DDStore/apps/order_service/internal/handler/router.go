package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupHttp(db *pgxpool.Pool) {
	productBaseUrl := "http://product_service:8888"
	userBaseUrl := "http://user_service:6666"

	orderService := Init(db, productBaseUrl)
	paymentService := InitPayment(db, userBaseUrl, productBaseUrl)

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

	payment := r.Group("/payment")
	{
		payment.POST("/:id", paymentService.CreatePaymentInfo())
		payment.POST("/checkout/:id", paymentService.CreateCheckoutSession())
		payment.POST("/webhook", paymentService.StripeWebhook())
	}

	r.Run(":7777")
}
