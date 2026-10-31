package handler

import (
	"context"
	"net/http"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/logger"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/metrics"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/repos"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func SetupHttp(db *pgxpool.Pool, rdb *redis.Client) {
	productBaseUrl := "http://product_service:8888"
	userBaseUrl := "http://user_service:6666"
	stockCache := repos.NewRedisStockCache(rdb)

	orderHandler := Init(db, productBaseUrl, stockCache)
	paymentHandler := InitPayment(db, userBaseUrl, productBaseUrl, stockCache)
	stockHandler := InitStock(db, productBaseUrl)

	// Cancels orders left unpaid for model.OrderPaymentTimeout and releases their stock.
	go orderHandler.orderService.RunExpiredOrderCanceller(context.Background(), model.ExpiredOrderSweepInterval)

	// gin.New instead of gin.Default: logger.Middleware replaces gin's text request log.
	r := gin.New()
	r.Use(logger.Middleware(), gin.Recovery())
	metrics.Serve()
	r.Use(metrics.Middleware())

	order := r.Group("/order")
	{
		order.POST("/:id", orderHandler.CreateOrder())
		order.POST("/guest", paymentHandler.GuestCheckout())
		order.GET("/:id", orderHandler.GetOrderById())
		order.GET("", orderHandler.GetOrders())
		order.GET("/history/:id", orderHandler.GetHistoryOrders())
		order.DELETE("/:id", orderHandler.CancelOrder())
		order.PUT("/:id", orderHandler.UpdateOrderState())
		order.POST("/completeOrder/:id", orderHandler.CompleteOrder())
		order.GET("/healthcheck", func(ctx *gin.Context) {
			ctx.JSON(http.StatusAccepted, gin.H{
				"message": "order service is healthy",
			})
		})
	}

	payment := r.Group("/payment")
	{
		payment.POST("/:id", paymentHandler.CreatePaymentInfo())
		payment.POST("/checkout/:id", paymentHandler.CreateCheckoutSession())
		payment.POST("/webhook", paymentHandler.StripeWebhook())
	}

	stock := r.Group("/stock")
	{
		stock.GET("", stockHandler.GetStocks())
		stock.POST("/:productId/add", stockHandler.AddStock())
	}

	r.Run(":7777")
}
