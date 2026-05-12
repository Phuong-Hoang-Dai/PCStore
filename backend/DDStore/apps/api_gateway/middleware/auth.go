package middleware

import (
	"net/http"
	"strconv"

	Err "github.com/Phuong-Hoang-Dai/DDStore/api_gateway/constant"
	"github.com/gin-gonic/gin"
)

func RequireRole(rqRole ...string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists {
			ctx.JSON(http.StatusForbidden, gin.H{"error": Err.ErrMissingRole})
			ctx.Abort()
			return
		}

		userRole := role.(string)
		for _, v := range rqRole {
			if v == userRole {
				ctx.Next()
				return
			}
		}

		ctx.JSON(http.StatusForbidden, gin.H{"error": Err.ErrNotAllowToAccess})
		ctx.Abort()
	}
}

func CheckIdForRoleUser() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists {
			responeError(http.StatusForbidden, Err.ErrMissingRole, ctx)
			ctx.Abort()
			return
		}
		if role == "user" {
			tokenId, exists := ctx.Get("id")
			if !exists {
				responeError(http.StatusForbidden, Err.ErrMissingId, ctx)
				ctx.Abort()
				return
			}
			id, err := strconv.Atoi(ctx.Param("id"))
			if err != nil {
				responeError(http.StatusBadRequest, err, ctx)
				ctx.Abort()
				return
			}
			if id != tokenId {
				responeError(http.StatusForbidden, Err.ErrNotAllowToAccess, ctx)
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
