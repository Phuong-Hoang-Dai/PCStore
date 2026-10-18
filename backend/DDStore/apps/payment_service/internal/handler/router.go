package handler

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func SetupHttp(db *gorm.DB) {
	_ = Init(db)

	r := gin.Default()

	payment := r.Group("/payment")
	{
		payment.GET("/healthcheck", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "payment service is healthy",
			})
		})
	}

	r.Run(":5555")
}
