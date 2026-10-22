package service

import (
	"context"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
)

type PaymentService struct {
	Repo PaymentRepository
}

func (p PaymentService) CreatePaymentInfo(orderId int, userId int, ctx context.Context) (status int, err error) {
	id, err := p.Repo.CreateOrderPayment(model.OrderPayment{
		OrderId: orderId,
		UserId:  userId,
	}, ctx)

	return id, err
}

func (p PaymentService) CreateCheckoutSession(orderId int, ctx context.Context) (url string, err error) {
	order, err := p.Repo.GetOrderPayment(orderId, ctx)
	if err != nil {
		return "", err
	}
}

func GetOrCreateCustomer(userId int, ctx context.Context) (int, error) {

}
