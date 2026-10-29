package model

import (
	"errors"
	"time"
)

var (
	ErrWrongFormat         = errors.New("has a wrong format")
	ErrOrderIsCreated      = errors.New("order is created")
	ErrOrderIsCanceled     = errors.New("order is canceled")
	ErrCannotGetUserId     = errors.New("can't get user id")
	ErrForbiddened         = errors.New("can't get someone else's history")
	ErrOutOfStock          = errors.New("be out of stock")
	ErrNotEnoughStock      = errors.New("be not enough stock")
	ErrOrderNotPaid        = errors.New("order is already paid or closed")
	ErrOrderNotCancellable = errors.New("only unpaid orders can be cancelled")
	ErrOrderPaymentExpired = errors.New("the payment window for this order has closed")
	ErrInvalidQuantity     = errors.New("quantity must be greater than 0")
	ErrProductNotFound     = errors.New("product not found")
)

const MaxLimit = 50

const (
	// OrderPaymentTimeout is how long an order may stay unpaid before it is
	// cancelled and its stock released.
	OrderPaymentTimeout = 3 * time.Hour
	// ExpiredOrderSweepInterval is how often unpaid orders are checked.
	ExpiredOrderSweepInterval = 5 * time.Minute
	// CheckoutExpiryMargin makes every Stripe session expire this long before the order
	// does, so a payment can't complete after the order was cancelled.
	CheckoutExpiryMargin = 10 * time.Minute
)
