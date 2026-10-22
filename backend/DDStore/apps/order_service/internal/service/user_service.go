package service

import (
	"context"
)

type UserService interface {
	GetUserById(userId int, ctx context.Context) (UserDTO, error)
}
