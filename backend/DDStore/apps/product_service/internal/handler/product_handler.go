package handler

import (
	"net/http"
	"strconv"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/repos"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductService struct {
	productRepos service.ProductRepos
}

func Init(db *gorm.DB) ProductService {
	return ProductService{productRepos: repos.NewMysqlProductRepo(db)}
}

func (p ProductService) CreateProduct() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var data model.Product

		if err := ctx.ShouldBind(&data); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}
		if data.CategoryID == 0 {
			data.CategoryID = 1
		}

		id, err := service.CreateProduct(data, p.productRepos)
		if err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Create Product succesfully",
			"data": gin.H{
				"id": id,
			},
		})
	}
}

func (p ProductService) GetProductById() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		var data model.Product
		if data, err = service.GetProductById(id, p.productRepos); err != nil {
			if err == gorm.ErrRecordNotFound {
				responeError(http.StatusNotFound, err, ctx)
				return
			} else {
				responeError(http.StatusInternalServerError, err, ctx)
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Retrieve Product succesfully",
			"data":    data,
		})
	}
}

func (p ProductService) UpdateProduct() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		data := model.Product{}
		if err := ctx.ShouldBind(&data); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
		}
		data.Id = id

		if err := service.UpdateProduct(data, p.productRepos); err != nil {
			if err == gorm.ErrRecordNotFound {
				responeError(http.StatusNotFound, err, ctx)
			} else {
				responeError(http.StatusInternalServerError, err, ctx)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Update Product succesfully",
			"data":    data,
		})
	}
}

func (p ProductService) DeleteProduct() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		if err := service.DeleteProduct(id, p.productRepos); err != nil {
			if err == gorm.ErrRecordNotFound {
				responeError(http.StatusNotFound, err, ctx)
			} else {
				responeError(http.StatusInternalServerError, err, ctx)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Delete Product succesfully",
			"data": gin.H{
				"id": id,
			},
		})
	}
}

func (p ProductService) GetProducts() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var paging model.Paging
		var err error

		paging.Limit, err = strconv.Atoi(ctx.DefaultQuery("limit", "10"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}
		paging.Offset, err = strconv.Atoi(ctx.DefaultQuery("offset", "0"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		var data []model.Product
		if data, err = service.GetProducts(&paging, p.productRepos); err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "List of Products retrieved successfully",
			"data":    data,
			"pagination": gin.H{
				"offset": paging.Offset,
				"limit":  paging.Limit,
				"total":  len(data),
			},
		})
	}
}

func (p ProductService) GetProductsByCate() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var paging model.Paging
		var err error
		cate := model.Category{}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}
		cate.Id = id

		paging.Limit, err = strconv.Atoi(ctx.DefaultQuery("limit", "10"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}
		paging.Offset, err = strconv.Atoi(ctx.DefaultQuery("offset", "0"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		var data []model.Product
		if data, err = service.GetProductsByCate(&paging, p.productRepos, cate); err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "List of Products by Category retrieved successfully",
			"data":    data,
			"pagination": gin.H{
				"offset": paging.Offset,
				"limit":  paging.Limit,
				"total":  len(data),
			},
		})
	}
}

func (p ProductService) GetStock() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var data []service.OrderItemsDto
		if err := ctx.ShouldBind(&data); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		if err := service.GetStock(data, p.productRepos); err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				responeError(http.StatusNotFound, err, ctx)
			case model.ErrOutOfStock:
				responeError(http.StatusBadRequest, err, ctx)
			default:
				responeError(http.StatusInternalServerError, err, ctx)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Update Stock successfully",
		})
	}
}

func (p ProductService) RestoreStock() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var data []service.OrderItemsDto
		if err := ctx.ShouldBind(&data); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		if err := service.RestoreStock(data, p.productRepos); err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				responeError(http.StatusNotFound, err, ctx)
			case model.ErrOutOfStock:
				responeError(http.StatusBadRequest, err, ctx)
			default:
				responeError(http.StatusInternalServerError, err, ctx)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Restore Stock successfully",
		})
	}
}

func (p ProductService) GetPriceProduct() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var data []service.OrderItemsDto
		if err := ctx.ShouldBind(&data); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		if err := service.GetPriceProduct(&data, p.productRepos); err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				responeError(http.StatusNotFound, err, ctx)
			case model.ErrOutOfStock:
				responeError(http.StatusBadRequest, err, ctx)
			default:
				responeError(http.StatusInternalServerError, err, ctx)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Prices of Products retrieved successfully",
			"data":    data,
		})
	}
}

func responeError(errCode int, err error, ctx *gin.Context) {
	ctx.JSON(errCode, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}
