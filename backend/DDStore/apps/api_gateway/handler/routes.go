package handler

import (
	"context"

	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/metrics"
	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/middleware"

	"github.com/gin-gonic/gin"
)

func SetupHttp() {
	userAPIUrl := "http://user_service:6666/user"
	productAPIUrl := "http://product_service:8888/product"
	cateAPIUrl := "http://product_service:8888/category"
	orderAPIUrl := "http://order_service:7777/order"
	paymentAPIUrl := "http://order_service:7777/payment"
	stockAPIUrl := "http://order_service:7777/stock"

	r := gin.Default()
	SetupRabbit(context.Background())
	metrics.Serve()
	r.Use(metrics.Middleware())
	r.Use(middleware.CORSMiddleware())
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			user := v1.Group("/user")
			{
				user.POST("/login", login(userAPIUrl))
				user.GET("/logout", logout())
				user.GET("/me", middleware.VerifyToken(), GetUserByToken(userAPIUrl))
				user.GET("/refresh", RefreshToken())
				user.POST("", createData(userAPIUrl))
				user.PUT("/:id", middleware.VerifyToken(), middleware.RequireRole("user", "admin", "owner"),
					middleware.CheckIdForRoleUser(), updateDataById(userAPIUrl))
				user.GET("", getDatas(userAPIUrl))
				user.GET("/:id", getDataById(userAPIUrl))
				user.DELETE("/:id", middleware.VerifyToken(), middleware.RequireRole("user", "admin", "owner"),
					middleware.CheckIdForRoleUser(), deleteDataById(userAPIUrl))
			}
			product := v1.Group("/product")
			{
				product.POST("", middleware.VerifyToken(), middleware.RequireRole("admin", "owner"), publishProduct())
				product.PUT("/:id", middleware.VerifyToken(), middleware.RequireRole("admin", "owner"), updateDataById(productAPIUrl))
				product.GET("", getDatas(productAPIUrl))
				product.GET("/:id", getDataById(productAPIUrl))
				product.GET("/category/:id", getProductsByCate(productAPIUrl+"/cate"))
				product.DELETE("/:id", middleware.VerifyToken(), middleware.RequireRole("owner"), deleteDataById(productAPIUrl))
			}
			cate := v1.Group("/category")
			{
				cate.POST("", middleware.VerifyToken(), middleware.RequireRole("admin", "owner"), createData(cateAPIUrl))
				cate.PUT("/:id", middleware.VerifyToken(), middleware.RequireRole("admin", "owner"), updateDataById(cateAPIUrl))
				cate.GET("", getDatas(cateAPIUrl))
				cate.GET("/:id", getDataById(cateAPIUrl))
				cate.DELETE("/:id", middleware.VerifyToken(), middleware.RequireRole("owner"), deleteDataById(cateAPIUrl))
			}
			order := v1.Group("/order", middleware.VerifyToken())
			{
				order.POST("", createOrder(orderAPIUrl))
				order.POST("/completeOrder/:id", middleware.RequireRole("admin", "owner"), completeOrderById(orderAPIUrl+"/completeOrder"))
				order.GET("/:id", middleware.RequireRole("user", "admin", "owner"), getOrderById(orderAPIUrl))
				order.GET("/history", middleware.RequireRole("user", "admin", "owner"), getHistoryOrderById(orderAPIUrl+"/history"))
				order.GET("", middleware.RequireRole("admin", "owner"), getDatas(orderAPIUrl))
				order.PUT("/:id", middleware.RequireRole("admin", "owner"), updateDataById(orderAPIUrl))
				order.DELETE("/:id", middleware.RequireRole("admin", "owner"), deleteDataById(orderAPIUrl))
			}
			payment := v1.Group("/payment", middleware.VerifyToken())
			{
				// completeOrderById forwards POST {baseUrl}/:id with an empty body, which is all checkout needs.
				payment.POST("/checkout/:id", middleware.RequireRole("user", "admin", "owner"), completeOrderById(paymentAPIUrl+"/checkout"))
			}
			stock := v1.Group("/stock", middleware.VerifyToken(), middleware.RequireRole("admin", "owner"))
			{
				stock.GET("", getStocks(stockAPIUrl))
				stock.POST("/:id/add", addStock(stockAPIUrl))
			}
		}
	}
	r.Run(":10000")
}
