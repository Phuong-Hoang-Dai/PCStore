package model

import "errors"

var (
	ErrWrongFormat     = errors.New("has a wrong format")
	ErrOrderIsCreated  = errors.New("order is created")
	ErrOrderIsCanceled = errors.New("order is canceled")
	ErrCannotGetUserId = errors.New("can't get user id")
	ErrForbiddened     = errors.New("can't get someone else's history")
	ErrOutOfStock      = errors.New("be out of stock")
)

const MaxLimit = 50
