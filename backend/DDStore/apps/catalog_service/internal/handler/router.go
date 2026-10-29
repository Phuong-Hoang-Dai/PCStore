package handler

import (
	"net/http"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/metrics"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupHttp(productService service.ProductService, cateService service.CateService) {
	productHandler := NewProductHandler(productService)
	cateHandler := NewCateHandler(cateService)
	r := gin.Default()
	metrics.Serve()
	r.Use(metrics.Middleware())
	r.Use(cors.Default())
	product := r.Group("/product")
	{
		product.POST("", productHandler.CreateProduct())
		product.POST("/ids", productHandler.GetProductsByIds())
		product.POST("/getprices", productHandler.GetPrices())
		product.GET("", productHandler.GetProducts())
		product.GET("/:id", productHandler.GetProductById())
		product.GET("/cate/:id", productHandler.GetProductsByCate())
		product.PUT("/:id", productHandler.UpdateProduct())
		product.PUT("/:id/state", productHandler.UpdateProductState())
		product.DELETE("/:id", productHandler.DeleteProduct())
		product.GET("/healthcheck", func(ctx *gin.Context) {
			ctx.JSON(http.StatusAccepted, gin.H{
				"message": "catalog service is healthy",
			})
		})
	}
	cate := r.Group("/category")
	{
		cate.POST("", cateHandler.CreateCate())
		cate.PUT("/:id", cateHandler.UpdateCate())
		cate.GET("", cateHandler.GetCates())
		cate.GET("/:id", cateHandler.GetCateById())
		cate.DELETE("/:id", cateHandler.DeleteCate())
	}
	r.Run(":8888")
}
