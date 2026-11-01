package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	Err "github.com/Phuong-Hoang-Dai/DDStore/api_gateway/constant"
	"github.com/gin-gonic/gin"
)

// getOrderById lets admins read any order but users only their own. The owner is only
// known once the order is fetched, so the check can't be a middleware on the URL id.
func getOrderById() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		target := orderService.String() + "/order/" + url.PathEscape(ctx.Param("id"))
		statusCode, body, err := fetch(ctx, http.MethodGet, target, nil)
		if err != nil {
			responeError(http.StatusBadGateway, err, ctx)
			return
		}

		if statusCode == http.StatusOK && ctx.GetString("role") == "user" {
			var resp struct {
				Data struct {
					UserId int `json:"userId"`
				} `json:"data"`
			}
			if err := json.Unmarshal(body, &resp); err != nil || resp.Data.UserId != ctx.GetInt("id") {
				responeError(http.StatusForbidden, Err.ErrNotAllowToAccess, ctx)
				return
			}
		}

		ctx.Data(statusCode, "application/json; charset=utf-8", body)
	}
}

// createOrder places the order for the logged-in user.
func createOrder() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		forward(ctx, orderService, fmt.Sprintf("/order/%d", ctx.GetInt("id")), ctx.Request.URL.RawQuery)
	}
}

// getHistoryOrderById lists the logged-in user's orders.
func getHistoryOrderById() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		forward(ctx, orderService, fmt.Sprintf("/order/history/%d", ctx.GetInt("id")), pagingQuery(ctx))
	}
}
