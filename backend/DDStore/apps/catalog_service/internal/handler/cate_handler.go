package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/service"
	"github.com/gin-gonic/gin"
)

type CateHandler struct {
	cateService service.CateService
}

func NewCateHandler(cateService service.CateService) CateHandler {
	return CateHandler{cateService: cateService}
}

func (c CateHandler) CreateCate() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var data model.Category

		if err := ctx.ShouldBind(&data); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		id, err := c.cateService.CreateCate(ctx.Request.Context(), data)
		if err != nil {
			responeCateError(err, ctx)
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

func (c CateHandler) GetCateById() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		data, err := c.cateService.GetCateById(ctx.Request.Context(), id)
		if err != nil {
			responeCateError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Retrieve Category succesfully",
			"data":    data,
		})
	}
}

func (c CateHandler) UpdateCate() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		data := model.Category{}
		if err := ctx.ShouldBind(&data); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}
		data.Id = id

		if err := c.cateService.UpdateCate(ctx.Request.Context(), data); err != nil {
			responeCateError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Update Category succesfully",
			"data":    data,
		})
	}
}

func (c CateHandler) DeleteCate() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		if err := c.cateService.DeleteCate(ctx.Request.Context(), id); err != nil {
			responeCateError(err, ctx)
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

func (c CateHandler) GetCates() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		data, err := c.cateService.GetCates(ctx.Request.Context())
		if err != nil {
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

func responeCateError(err error, ctx *gin.Context) {
	switch {
	case errors.Is(err, model.ErrNotFound):
		responeError(http.StatusNotFound, err, ctx)
	case errors.Is(err, model.ErrCategoryIsInvalid):
		responeError(http.StatusBadRequest, err, ctx)
	case errors.Is(err, model.ErrCategoryInUse):
		responeError(http.StatusConflict, err, ctx)
	default:
		responeError(http.StatusInternalServerError, err, ctx)
	}
}
