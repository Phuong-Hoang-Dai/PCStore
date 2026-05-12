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

type CateService struct {
	cateRepos service.CateRepos
}

func InitCateService(db *gorm.DB) CateService {
	return CateService{cateRepos: repos.NewMysqlCateRepo(db)}
}

func (c CateService) CreateCate() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var data model.Category

		if err := ctx.ShouldBind(&data); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		id, err := service.CreateCate(data, c.cateRepos)
		if err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Create Category succesfully",
			"data": gin.H{
				"id": id,
			},
		})
	}
}

func (c CateService) GetCateById() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		var data model.Category
		if data, err = service.GetCateById(id, c.cateRepos); err != nil {
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
			"message": "Retrieve Category succesfully",
			"data":    data,
		})
	}
}

func (c CateService) UpdateCate() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		data := model.Category{}
		if err := ctx.ShouldBind(&data); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
		}
		data.Id = id

		if err := service.UpdateCate(data, c.cateRepos); err != nil {
			if err == gorm.ErrRecordNotFound {
				responeError(http.StatusNotFound, err, ctx)
			} else {
				responeError(http.StatusInternalServerError, err, ctx)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Update Category succesfully",
			"data":    data,
		})
	}
}

func (c CateService) DeleteCate() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		if err := service.DeleteCate(id, c.cateRepos); err != nil {
			if err == gorm.ErrRecordNotFound {
				responeError(http.StatusNotFound, err, ctx)
			} else {
				responeError(http.StatusInternalServerError, err, ctx)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Delete Category succesfully",
			"data": gin.H{
				"id": id,
			},
		})
	}
}

func (c CateService) GetCates() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var err error

		var data []model.Category
		if data, err = service.GetCates(c.cateRepos); err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "List of Categories retrieved successfully",
			"data":    data,
		})
	}
}
