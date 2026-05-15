package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/repos"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(db *mongo.Client, redis *redis.Client) ProductHandler {
	col := db.Database(configs.Cfg.DBName).Collection("products")
	repo := repos.NewMongoProductRepo(col)
	return ProductHandler{productService: service.NewProductService(repo, redis)}
}

func (p ProductHandler) CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.Product

		if err := c.ShouldBind(&data); err != nil {
			responeError(http.StatusBadRequest, err, c)
			return
		}

		id, err := p.productService.CreateProduct(c.Request.Context(), data)
		if err != nil {
			responeError(http.StatusInternalServerError, err, c)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Create Product succesfully",
			"data":    gin.H{"id": id.Hex()},
		})
	}
}

func (p ProductHandler) GetProductById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := bson.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, c)
			return
		}

		data, err := p.productService.GetProductById(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				responeError(http.StatusNotFound, err, c)
			} else {
				responeError(http.StatusInternalServerError, err, c)
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Retrieve Product succesfully",
			"data":    data,
		})
	}
}

func (p ProductHandler) UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := bson.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, c)
			return
		}

		var data model.Product
		if err := c.ShouldBind(&data); err != nil {
			responeError(http.StatusBadRequest, err, c)
			return
		}
		data.Id = id

		if err := p.productService.UpdateProduct(c.Request.Context(), data); err != nil {
			if errors.Is(err, model.ErrNotFound) {
				responeError(http.StatusNotFound, err, c)
			} else {
				responeError(http.StatusInternalServerError, err, c)
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Update Product succesfully",
			"data":    data,
		})
	}
}

func (p ProductHandler) DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := bson.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, c)
			return
		}

		if err := p.productService.DeleteProduct(c.Request.Context(), id); err != nil {
			if errors.Is(err, model.ErrNotFound) {
				responeError(http.StatusNotFound, err, c)
			} else {
				responeError(http.StatusInternalServerError, err, c)
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Delete Product succesfully",
			"data":    gin.H{"id": id.Hex()},
		})
	}
}

func (p ProductHandler) GetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var paging model.Paging
		var err error

		paging.Limit, err = strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil {
			responeError(http.StatusBadRequest, err, c)
			return
		}
		paging.Offset, err = strconv.Atoi(c.DefaultQuery("offset", "0"))
		if err != nil {
			responeError(http.StatusBadRequest, err, c)
			return
		}

		data, err := p.productService.GetProducts(c.Request.Context(), &paging)
		if err != nil {
			responeError(http.StatusInternalServerError, err, c)
			return
		}

		c.JSON(http.StatusOK, gin.H{
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

func (p ProductHandler) GetProductsByCate() gin.HandlerFunc {
	return func(c *gin.Context) {
		cateID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, c)
			return
		}

		var paging model.Paging
		paging.Limit, err = strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil {
			responeError(http.StatusBadRequest, err, c)
			return
		}
		paging.Offset, err = strconv.Atoi(c.DefaultQuery("offset", "0"))
		if err != nil {
			responeError(http.StatusBadRequest, err, c)
			return
		}

		cate := model.Category{Id: cateID}
		data, err := p.productService.GetProductsByCate(c.Request.Context(), &paging, cate)
		if err != nil {
			responeError(http.StatusInternalServerError, err, c)
			return
		}

		c.JSON(http.StatusOK, gin.H{
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

func responeError(errCode int, err error, c *gin.Context) {
	c.JSON(errCode, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}
