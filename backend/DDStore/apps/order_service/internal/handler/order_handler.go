package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/repos"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/gorm"
)

type OrderHandler struct {
	orderService service.OrderService
}

func Init(db *pgxpool.Pool, url string) OrderHandler {
	o := OrderHandler{
		orderService: service.NewOrderService(
			repos.NewPgOrderRepo(db),
			productHTTPClient{baseURL: fmt.Sprint(url, "/product")},
			service.StockService{Repos: repos.NewPgStockRepo(db), ProductService: productHTTPClient{baseURL: fmt.Sprint(url, "/product")}},
			service.SequencerService{Repos: repos.NewPgSequencerkRepo(db)},
		),
	}

	return o
}

func (o OrderHandler) CreateOrder() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userId, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		bodyData := struct {
			UserId int                `json:"userId"`
			Data   []service.OrderDTO `json:"data"`
		}{}

		if err := ctx.ShouldBind(&bodyData); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}
		bodyData.UserId = userId
		if userId == 0 {
			responeError(http.StatusBadRequest, model.ErrCannotGetUserId, ctx)
			return
		}

		data := bodyData.Data
		id, err := o.orderService.CreateOrder(userId, data, ctx)

		if err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Create Order succesfully",
			"data": gin.H{
				"id": id,
			},
		})
	}
}

func (o OrderHandler) CancelOrder() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		err = o.orderService.CancelOrder(id, ctx)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			responeError(http.StatusNotFound, err, ctx)
			return
		}

		if err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Cancel Order succesfully",
		})
	}
}

func (o OrderHandler) GetOrderById() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		data, err := o.orderService.GetOrderById(id, ctx)

		if err == gorm.ErrRecordNotFound {
			responeError(http.StatusNotFound, err, ctx)
			return
		}

		if err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Retrieve Order succesfully",
			"data":    data,
		})
	}
}

func (o OrderHandler) UpdateOrderState() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		data := struct {
			State model.OrderState `json:"state"`
		}{}

		if err := ctx.ShouldBind(&data); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
		}

		if err := o.orderService.UpdateOrder(id, data.State, ctx); err != nil {
			if err == gorm.ErrRecordNotFound {
				responeError(http.StatusNotFound, err, ctx)
			} else {
				responeError(http.StatusInternalServerError, err, ctx)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Update State Order succesfully",
		})
	}
}

func (o OrderHandler) CompleteOrder() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		if err := o.orderService.UpdateOrder(id, model.Completed, ctx); err != nil {
			if err == gorm.ErrRecordNotFound {
				responeError(http.StatusNotFound, err, ctx)
			} else {
				responeError(http.StatusInternalServerError, err, ctx)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Completed Order succesfully",
		})
	}
}

func (o OrderHandler) GetOrders() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var p model.Paging
		var err error

		p.Limit, err = strconv.Atoi(ctx.DefaultQuery("limit", "0"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}
		p.Offset, err = strconv.Atoi(ctx.DefaultQuery("offset", "0"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		var data []model.Order
		if data, err = o.orderService.GetOrders(&p, ctx); err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "List of Orders retrieved successfully",
			"data":    data,
			"pagination": gin.H{
				"offset": p.Offset,
				"limit":  p.Limit,
				"total":  len(data),
			},
		})
	}
}

func (o OrderHandler) GetHistoryOrders() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var p model.Paging
		var err error

		userId, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		p.Limit, err = strconv.Atoi(ctx.DefaultQuery("limit", "0"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}
		p.Offset, err = strconv.Atoi(ctx.DefaultQuery("offset", "0"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		var data []model.Order
		if data, err = o.orderService.GetHistoryOrders(userId, &p, ctx); err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "History Orders retrieved successfully",
			"data":    data,
			"pagination": gin.H{
				"offset": p.Offset,
				"limit":  p.Limit,
				"total":  len(data),
			},
		})
	}
}

func responeError(errCode int, err error, ctx *gin.Context) {
	ctx.JSON(errCode, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}
