package handler

import (
	"net/http"
	"strconv"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/repos"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
)

type ProductHandler struct {
	productService service.ProductService
}

func Init(db *mongo.Client) ProductHandler {
	col := db.Database("ddstore").Collection("product")
	repos := repos.NewMongoProductRepo(col)
	return ProductHandler{productService: service.NewProductService(repos)}
}

func (p ProductHandler) CreateProduct() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var data model.Product

		if err := ctx.ShouldBind(&data); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		id, err := p.productService.CreateProduct(data)
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

func (p ProductHandler) GetProductById() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		var data model.Product
		if data, err = p.productService.GetProductById(id); err != nil {
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

func (p ProductHandler) UpdateProduct() func(ctx *gin.Context) {
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

		if err := p.productService.UpdateProduct(data); err != nil {
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

func (p ProductHandler) DeleteProduct() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		if err := p.productService.DeleteProduct(id); err != nil {
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

func (p ProductHandler) GetProducts() func(ctx *gin.Context) {
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
		if data, err = p.productService.GetProducts(&paging); err != nil {
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

func (p ProductHandler) GetProductsByCate() func(ctx *gin.Context) {
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
		if data, err = p.productService.GetProductsByCate(&paging, cate); err != nil {
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

func responeError(errCode int, err error, ctx *gin.Context) {
	ctx.JSON(errCode, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}
