package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupHttp(db *mongo.Client, redis *redis.Client) {
	productHandler := NewProductHandler(db, redis)
	//	cateService := InitCateService(db)
	r := gin.Default()
	r.Use(cors.Default())
	product := r.Group("/product")
	{
		product.POST("", productHandler.CreateProduct())
		product.GET("", productHandler.GetProducts())
		product.GET("/:id", productHandler.GetProductById())
		product.GET("/cate/:id", productHandler.GetProductsByCate())
		product.PUT("/:id", productHandler.UpdateProduct())
		product.DELETE("/:id", productHandler.DeleteProduct())
	}
	// cate := r.Group("/category")
	// {
	// 	cate.POST("", cateService.CreateCate())
	// 	cate.PUT("/:id", cateService.UpdateCate())
	// 	cate.GET("", cateService.GetCates())
	// 	cate.GET("/:id", cateService.GetCateById())
	// 	cate.DELETE("/:id", cateService.DeleteCate())
	// }
	r.Run(":8888")
}
