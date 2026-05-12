package handler

import (
	"net/http"
	"strconv"

	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/repos"
	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService struct {
	userRepos service.UserRepos
}

func Init(db *gorm.DB) UserService {
	return UserService{
		userRepos: repos.NewMysqlUserRepo(db),
	}
}

func (u UserService) CreateUser() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var data service.UserCreateDTO

		if err := ctx.ShouldBind(&data); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		id, err := service.CreateUser(data, u.userRepos)
		if err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Create User succesfully",
			"data": gin.H{
				"id": id,
			},
		})
	}
}

func (u UserService) GetUserById() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		var data service.UserResponeDTO
		data, err = service.GetUserById(id, u.userRepos)
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
			"message": "Retrieve User succesfully",
			"data":    data,
		})
	}
}

func (u UserService) UpdateUser() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		data := service.UserUpdateDTO{}
		if err := ctx.ShouldBind(&data); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		data.Id = id
		if err := service.UpdateUser(data, u.userRepos); err != nil {
			if err == gorm.ErrRecordNotFound {
				responeError(http.StatusNotFound, err, ctx)
			} else {
				responeError(http.StatusInternalServerError, err, ctx)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Update User succesfully",
			"data":    data,
		})
	}
}

func (u UserService) DeleteUser() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		if err := service.DeleteUser(id, u.userRepos); err != nil {
			if err == gorm.ErrRecordNotFound {
				responeError(http.StatusNotFound, err, ctx)
			} else {
				responeError(http.StatusInternalServerError, err, ctx)
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Delete User succesfully",
			"data": gin.H{
				"id": id,
			},
		})
	}
}

func (u UserService) GetUsers() func(ctx *gin.Context) {
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

		var data []service.UserResponeDTO

		if data, err = service.GetUsers(&p, u.userRepos); err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "List of Users retrieved successfully",
			"data":    data,
			"pagination": gin.H{
				"offset": p.Offset,
				"limit":  p.Limit,
				"total":  len(data),
			},
		})
	}
}

func (u UserService) VerifyPassword() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var data service.UserDTO
		if err := ctx.ShouldBind(&data); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}
		if user, err := service.VerifyPassword(data, u.userRepos); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "correct password successfully",
				"data":    user,
			})
		}
	}
}

func responeError(errCode int, err error, ctx *gin.Context) {
	ctx.JSON(errCode, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}
