package service

import (
	"context"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
)

type OrderRepository interface {
	CreateOrder(data *model.Order, ctx context.Context) (int, error)
	UpdateOrder(data model.Order, ctx context.Context) error
	GetOrderById(id int, ctx context.Context) (model.Order, error)
	GetOrders(p model.Paging, ctx context.Context) ([]model.Order, error)
	GetHistoryOrders(id int, p model.Paging, ctx context.Context) ([]model.Order, error)
	DeleteOrder(id int, ctx context.Context) error
}

type StockRepository interface {
	ReserveStock(productId string, qnantity int, orderId int, ctx context.Context) error
}

type SequenecerRepository interface {
	GetOrderId(ctx context.Context) (int, error)
}
