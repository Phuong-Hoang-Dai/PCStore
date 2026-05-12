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

func login(baseUrl string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			responeError(http.StatusBadRequest, Err.ErrCannotReadBody, ctx)
			return
		}

		bodyData, statusCode, err := makeRequest(http.MethodPost, fmt.Sprint(baseUrl, "/verifyPassword"), bytes.NewBuffer(body))
		if err != nil {
			responeError(statusCode, err, ctx)
			return
		}

		if statusCode != http.StatusOK {
			ctx.JSON(statusCode, bodyData)
			return
		}

		data := model.UseResp{}
		if err := json.Unmarshal(bodyData["data"], &data); err != nil {
			responeError(http.StatusInternalServerError, Err.ErrCannotReadBody, ctx)
			return
		}

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

func GetUserByToken(baseUrl string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetInt("id")

		bodyData, statusCode, err := makeRequest(http.MethodGet, fmt.Sprintf("%s/%d", baseUrl, userId), bytes.NewBuffer([]byte{}))
		if err != nil {
			responeError(statusCode, err, ctx)
			return
		}

		ctx.JSON(statusCode, bodyData)
	}
}
