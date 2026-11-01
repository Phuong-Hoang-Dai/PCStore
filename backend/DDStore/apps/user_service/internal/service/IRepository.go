package service

import (
	"context"

	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/model"
)

type UserRepos interface {
	// CreateUser registers a user. A guest account (no password) with the same email is
	// claimed instead of creating a second account; an email that already belongs to a
	// registered account returns model.ErrEmailAlreadyRegistered.
	CreateUser(data *model.User, ctx context.Context) (int, error)
	// FindOrCreateGuest returns the account using data.Email, or creates one without a password.
	FindOrCreateGuest(data *model.User, ctx context.Context) (int, error)
	UpdateUser(data model.User, ctx context.Context) error
	GetUserById(id int, data *model.User, ctx context.Context) error
	GetUserByName(name string, data *model.User, ctx context.Context) error
	GetUsers(p model.Paging, data *[]model.User, ctx context.Context) error
	DeleteUser(id int, ctx context.Context) error
}
