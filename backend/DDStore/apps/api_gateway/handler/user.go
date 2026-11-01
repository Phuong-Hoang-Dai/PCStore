package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/configs"
	Err "github.com/Phuong-Hoang-Dai/DDStore/api_gateway/constant"
	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/model"
	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/utils"
	"github.com/gin-gonic/gin"
)

// signUp forwards a registration. user_service treats a user without a password as a guest
// (an account for orders placed without logging in), which only order_service may create, so
// a public sign-up must have one.
func signUp() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			responeError(http.StatusBadRequest, Err.ErrCannotReadBody, ctx)
			return
		}

		var data struct {
			Password string `json:"password"`
		}
		if err := json.Unmarshal(body, &data); err != nil || data.Password == "" {
			responeError(http.StatusBadRequest, Err.ErrPasswordRequired, ctx)
			return
		}

		// Read above, so it is put back for the service.
		ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
		forwardTo(userService)(ctx)
	}
}

// login checks the password with user_service and, if it is right, sets the token cookies.
func login() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		statusCode, body, err := fetch(ctx, http.MethodPost, userService.String()+"/user/verifyPassword", ctx.Request.Body)
		if err != nil {
			responeError(http.StatusBadGateway, err, ctx)
			return
		}

		if statusCode != http.StatusOK {
			ctx.Data(statusCode, "application/json; charset=utf-8", body)
			return
		}

		var resp struct {
			Data model.UseResp `json:"data"`
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			responeError(http.StatusInternalServerError, Err.ErrCannotReadBody, ctx)
			return
		}
		data := resp.Data

		var access_token, refresh_token string
		if access_token, err = utils.GenerateJwt(data, configs.Cfg.AccessTokenExpireIn); err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}
		if refresh_token, err = utils.GenerateJwt(data, configs.Cfg.RefreshTokenExpireIn); err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}

		access_token_expire, err := utils.ParseCustomDuration(configs.Cfg.AccessTokenExpireIn)
		if err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}
		refresh_token_expire, err := utils.ParseCustomDuration(configs.Cfg.RefreshTokenExpireIn)
		if err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}

		ctx.SetCookie(
			"access_token", access_token, int(access_token_expire.Seconds()), "/", "localhost", false, true,
		)

		ctx.SetCookie(
			"refresh_token", refresh_token, int(refresh_token_expire.Seconds()), "/api/v1/user/refresh",
			"localhost", false, true,
		)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Login successfully",
		})
	}
}

func RefreshToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr, err := ctx.Cookie("refresh_token")
		if err != nil {
			responeError(http.StatusBadRequest, Err.ErrCookieMissing, ctx)
			return
		}
		access_token, err := utils.RefreshToken(tokenStr)
		if err != nil {
			responeError(http.StatusInternalServerError, Err.ErrCookieMissing, ctx)
			return
		}

		access_token_expire, err := utils.ParseCustomDuration(configs.Cfg.AccessTokenExpireIn)
		if err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}
		ctx.SetCookie(
			"access_token", access_token, int(access_token_expire.Seconds()), "/",
			"localhost", false, true,
		)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Refresh token successfully",
		})
	}
}

func logout() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		ctx.SetCookie(
			"access_token", "", 0, "/", "localhost", false, true,
		)

		ctx.SetCookie(
			"refresh_token", "", 0, "/api/v1/user/refresh",
			"localhost", false, true,
		)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Login successfully",
		})
	}
}

// GetUserByToken forwards to the logged-in user's own record.
func GetUserByToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		forward(ctx, userService, fmt.Sprintf("/user/%d", ctx.GetInt("id")), "")
	}
}
