package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() func(ctx *gin.Context) {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Transaction-Id"},
		ExposeHeaders:    []string{"Content-Length", "X-Transaction-Id"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
