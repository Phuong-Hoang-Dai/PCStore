package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getProductsByCate(baseUrl string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}
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
		bodyData, statusCode, err := makeRequest(http.MethodGet,
			fmt.Sprintf("%s/%d?limit=%d&offset=%d", baseUrl, id, limit, offset),
			bytes.NewBuffer([]byte{}))
		if err != nil {
			responeError(statusCode, err, ctx)
			return
		}

		ctx.JSON(statusCode, bodyData)
	}
}
