package handler

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func SetupHttp(db *gorm.DB) {
	userService := Init(db)

	r := gin.Default()

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
