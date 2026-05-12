package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	Err "github.com/Phuong-Hoang-Dai/DDStore/api_gateway/constant"
	"github.com/gin-gonic/gin"
)

func completeOrderById(baseUrl string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}
		bodyData, statusCode, err := makeRequest(http.MethodPost, fmt.Sprintf("%s/%d", baseUrl, id), bytes.NewBuffer([]byte{}))
		if err != nil {
			responeError(statusCode, err, ctx)
			return
		}

		ctx.JSON(statusCode, bodyData)
	}
}

func createOrder(baseUrl string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			responeError(http.StatusBadRequest, Err.ErrCannotReadBody, ctx)
		}
		id := ctx.GetInt("id")
		bodyData, statusCode, err := makeRequest(http.MethodPost, fmt.Sprintf("%s/%d", baseUrl, id), bytes.NewBuffer(body))
		if err != nil {
			responeError(statusCode, err, ctx)
			return
		}

		ctx.JSON(statusCode, bodyData)
	}
}

func getHistoryOrderById(baseUrl string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.GetInt("id")
		limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "0"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}
		offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}
		bodyData, statusCode, err := makeRequest(http.MethodGet, fmt.Sprintf("%s/%d?limit=%d&offset=%d", baseUrl, id, limit, offset), bytes.NewBuffer([]byte{}))

		if err != nil {
			responeError(statusCode, err, ctx)
			return
		}

		ctx.JSON(statusCode, bodyData)
	}
}
