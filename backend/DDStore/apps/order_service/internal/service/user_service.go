package service

import (
	"context"
)

type UserService interface {
	GetUserById(userId int, ctx context.Context) (UserDTO, error)
	// GetOrCreateGuest returns the id of the account using the guest's email; the user service
	// creates one without a password if there is none.
	GetOrCreateGuest(guest GuestDTO, ctx context.Context) (int, error)
}
