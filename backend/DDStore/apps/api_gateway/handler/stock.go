package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	Err "github.com/Phuong-Hoang-Dai/DDStore/api_gateway/constant"
	"github.com/gin-gonic/gin"
)

// getStocks forwards GET /stock with its query string (?ids=a,b,c).
func getStocks(baseUrl string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		target := baseUrl
		if query := ctx.Request.URL.RawQuery; query != "" {
			target += "?" + query
		}

		bodyData, statusCode, err := makeRequest(http.MethodGet, target, bytes.NewBuffer([]byte{}))
		if err != nil {
			responeError(statusCode, err, ctx)
			return
		}

		ctx.JSON(statusCode, bodyData)
	}
}

// addStock forwards POST /stock/:id/add with its {"quantity": n} body.
func addStock(baseUrl string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			responeError(http.StatusBadRequest, Err.ErrCannotReadBody, ctx)
			return
		}

		target := fmt.Sprintf("%s/%s/add", baseUrl, url.PathEscape(ctx.Param("id")))
		bodyData, statusCode, err := makeRequest(http.MethodPost, target, bytes.NewBuffer(body))
		if err != nil {
			responeError(statusCode, err, ctx)
			return
		}

		ctx.JSON(statusCode, bodyData)
	}
}
