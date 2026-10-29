package handler

import (
	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/metrics"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupHttp(db *pgxpool.Pool) {
	userService := Init(db)

	r := gin.Default()
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
	}

	r.Run(":6666")
}
