package service

import (
	"context"

	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/model"
)

type UserRepos interface {
	CreateUser(data *model.User, ctx context.Context) (int, error)
	UpdateUser(data model.User, ctx context.Context) error
	GetUserById(id int, data *model.User, ctx context.Context) error
	GetUserByName(name string, data *model.User, ctx context.Context) error
	GetUsers(p model.Paging, data *[]model.User, ctx context.Context) error
	DeleteUser(id int, ctx context.Context) error
}
