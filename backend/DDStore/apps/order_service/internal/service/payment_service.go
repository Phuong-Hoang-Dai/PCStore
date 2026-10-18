package service

import (
	"context"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
)

type PaymentService interface {
	CreatePaymentInfo(order model.Order, orderItems []model.OrderItem, ctx context.Context) (status int, err error)
}
