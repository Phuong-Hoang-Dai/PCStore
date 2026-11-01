package handler

import (
	"net/http"

	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/logger"
	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/metrics"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupHttp(db *pgxpool.Pool) {
	userService := Init(db)

	// gin.New instead of gin.Default: logger.Middleware replaces gin's text request log.
	r := gin.New()
	r.Use(logger.Middleware(), gin.Recovery())
	metrics.Serve()
	r.Use(metrics.Middleware())

	user := r.Group("/user")
	{
		user.POST("", userService.CreateUser())
		user.PUT("/:id", userService.UpdateUser())
		user.GET("/:id", userService.GetUserById())
		user.GET("", userService.GetUsers())
		user.DELETE("/:id", userService.DeleteUser())
		user.POST("/verifyPassword", userService.VerifyPassword())
		user.GET("/healthcheck", func(ctx *gin.Context) {
			ctx.JSON(http.StatusAccepted, gin.H{
				"message": "user service is healthy",
			})
		})
	}

	r.Run(":6666")
}
