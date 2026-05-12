package service

import (
	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/model"
)

type UserRepos interface {
	CreateUser(data *model.User) (int, error)
	UpdateUser(data model.User) error
	GetUserById(id int, data *model.User) error
	GetUserByName(name string, data *model.User) error
	GetUsers(p model.Paging, data *[]model.User) error
	DeleteUser(id int) error
}
