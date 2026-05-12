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
	"gorm.io/gorm"
)

type OrderService struct {
	orderRepos     service.OrderRepository
	productService productHTTPClient
}

func Init(db *gorm.DB, url string) OrderService {
	o := OrderService{
		orderRepos:     repos.NewMysqlOrderRepo(db),
		productService: productHTTPClient{baseURL: fmt.Sprint(url, "/product")},
	}

	return o
}

func (o OrderService) CreateOrder() func(ctx *gin.Context) {
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
		id, err := service.CreateOrder(userId, data, o.orderRepos, o.productService)

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

func (o OrderService) CancelOrder() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		err = service.CancelOrder(id, o.orderRepos, o.productService)

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

func (o OrderService) GetOrderById() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		data, err := service.GetOrderById(id, o.orderRepos)

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

func (o OrderService) UpdateOrderState() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		data := struct {
			State int `json:"state"`
		}{}

		if err := ctx.ShouldBind(&data); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
		}

		if err := service.UpdateOrder(id, data.State, o.orderRepos); err != nil {
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

func (o OrderService) CompleteOrder() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		if err := service.UpdateOrder(id, model.Completed, o.orderRepos); err != nil {
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

func (o OrderService) GetOrders() func(ctx *gin.Context) {
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
		if data, err = service.GetOrders(&p, o.orderRepos); err != nil {
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

func (o OrderService) GetHistoryOrders() func(ctx *gin.Context) {
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
		if data, err = service.GetHistoryOrders(userId, &p, o.orderRepos); err != nil {
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
