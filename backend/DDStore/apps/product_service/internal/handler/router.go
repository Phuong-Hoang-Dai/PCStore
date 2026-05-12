package handler

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func SetupHttp(db *gorm.DB) {
	productService := Init(db)
	cateService := InitCateService(db)
	r := gin.Default()
	product := r.Group("/product")
	{
		product.POST("", productService.CreateProduct())
		product.POST("/getstock", productService.GetStock())
		product.POST("/restore", productService.RestoreStock())
		product.GET("", productService.GetProducts())
		product.POST("/getprice", productService.GetPriceProduct())
		product.GET("/:id", productService.GetProductById())
		product.GET("/cate/:id", productService.GetProductsByCate())
		product.PUT("/:id", productService.UpdateProduct())
		product.DELETE("/:id", productService.DeleteProduct())
	}
	cate := r.Group("/category")
	{
		cate.POST("", cateService.CreateCate())
		cate.PUT("/:id", cateService.UpdateCate())
		cate.GET("", cateService.GetCates())
		cate.GET("/:id", cateService.GetCateById())
		cate.DELETE("/:id", cateService.DeleteCate())

	}
	r.Run(":8888")
}
