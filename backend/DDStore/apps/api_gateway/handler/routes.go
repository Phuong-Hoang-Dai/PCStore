package handler

import (
	"context"
	"net/http"

	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/logger"
	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/metrics"
	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/middleware"

	"github.com/gin-gonic/gin"
)

// SetupHttp serves the API. Most routes check the caller's token and role, then forward the
// request unchanged to the service under the same path without /api/v1 (see proxy.go).
func SetupHttp() {
	// gin.New instead of gin.Default: logger.Middleware replaces gin's text request log.
	r := gin.New()
	r.Use(logger.Middleware(), gin.Recovery())
	SetupRabbit(context.Background())
	metrics.Serve()
	r.Use(metrics.Middleware())
	r.Use(middleware.CORSMiddleware())
	r.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusAccepted, gin.H{
			"message": "api gateway is healthy",
		})
	})
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			user := v1.Group("/user")
			{
				user.POST("/login", login())
				user.GET("/logout", logout())
				user.GET("/me", middleware.VerifyToken(), GetUserByToken())
				user.GET("/refresh", RefreshToken())
				user.POST("", signUp())
				user.PUT("/:id", middleware.VerifyToken(), middleware.RequireRole("user", "admin", "owner"),
					middleware.CheckIdForRoleUser(), forwardTo(userService))
				user.GET("", forwardList(userService))
				user.GET("/:id", forwardTo(userService))
				user.DELETE("/:id", middleware.VerifyToken(), middleware.RequireRole("user", "admin", "owner"),
					middleware.CheckIdForRoleUser(), forwardTo(userService))
			}
			product := v1.Group("/product")
			{
				product.POST("", middleware.VerifyToken(), middleware.RequireRole("admin", "owner"), publishProduct())
				product.PUT("/:id", middleware.VerifyToken(), middleware.RequireRole("admin", "owner"), forwardTo(productService))
				product.GET("", forwardList(productService))
				product.GET("/:id", forwardTo(productService))
				product.GET("/category/:id", getProductsByCate())
				product.DELETE("/:id", middleware.VerifyToken(), middleware.RequireRole("owner"), forwardTo(productService))
			}
			cate := v1.Group("/category")
			{
				cate.POST("", middleware.VerifyToken(), middleware.RequireRole("admin", "owner"), forwardTo(productService))
				cate.PUT("/:id", middleware.VerifyToken(), middleware.RequireRole("admin", "owner"), forwardTo(productService))
				cate.GET("", forwardTo(productService))
				cate.GET("/:id", forwardTo(productService))
				cate.DELETE("/:id", middleware.VerifyToken(), middleware.RequireRole("owner"), forwardTo(productService))
			}
			order := v1.Group("/order", middleware.VerifyToken())
			{
				order.POST("", createOrder())
				order.POST("/completeOrder/:id", middleware.RequireRole("admin", "owner"), forwardTo(orderService))
				order.GET("/:id", middleware.RequireRole("user", "admin", "owner"), getOrderById())
				order.GET("/history", middleware.RequireRole("user", "admin", "owner"), getHistoryOrderById())
				order.GET("", middleware.RequireRole("admin", "owner"), forwardList(orderService))
				order.PUT("/:id", middleware.RequireRole("admin", "owner"), forwardTo(orderService))
				order.DELETE("/:id", middleware.RequireRole("admin", "owner"), forwardTo(orderService))
			}
			// No login: order_service gets (or creates) the account for the guest's email and
			// returns the Stripe checkout URL along with the order id.
			v1.POST("/order/guest", forwardTo(orderService))
			payment := v1.Group("/payment", middleware.VerifyToken())
			{
				payment.POST("/checkout/:id", middleware.RequireRole("user", "admin", "owner"), forwardTo(orderService))
			}
			stock := v1.Group("/stock", middleware.VerifyToken(), middleware.RequireRole("admin", "owner"))
			{
				stock.GET("", forwardTo(orderService))
				stock.POST("/:id/add", forwardTo(orderService))
			}
		}
	}
	r.Run(":10000")
}
