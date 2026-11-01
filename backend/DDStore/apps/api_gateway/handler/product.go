package handler

import (
	"io"
	"net/http"

	Err "github.com/Phuong-Hoang-Dai/DDStore/api_gateway/constant"
	"github.com/gin-gonic/gin"
)

// publishProduct hands the new product to catalog_service through RabbitMQ.
func publishProduct() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			responeError(http.StatusBadRequest, Err.ErrCannotReadBody, ctx)
			return
		}
		RabbitHandler.PublishProduct(ctx, body)

		ctx.JSON(http.StatusAccepted, "")
	}
}

// getProductsByCate forwards to catalog_service, which serves it at /product/cate/:id.
func getProductsByCate() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		forward(ctx, productService, "/product/cate/"+ctx.Param("id"), pagingQuery(ctx))
	}
}
