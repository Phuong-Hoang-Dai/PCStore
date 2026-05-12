package middleware

import (
	"net/http"

	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/constant"
	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		tokenStr, err := ctx.Cookie("access_token")
		if err != nil {
			responeError(http.StatusBadRequest, constant.ErrCookieMissing, ctx)
			ctx.Abort()
			return
		}

		if token, err := utils.VerifyJwt(tokenStr); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			ctx.Abort()
		} else {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				userId := int(claims["sub"].(float64))
				ctx.Set("id", userId)
				ctx.Set("role", claims["role"])
			}

			ctx.Next()
		}
	}
}

func responeError(errCode int, err error, ctx *gin.Context) {
	ctx.JSON(errCode, gin.H{
		"error": err.Error(),
	})
}
