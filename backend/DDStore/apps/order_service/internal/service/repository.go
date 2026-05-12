package service

import (
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
)

type OrderRepository interface {
	CreateOrder(data *model.Order) (int, error)
	UpdateOrder(data model.Order) error
	GetOrderById(id int) (model.Order, error)
	GetOrders(p model.Paging) ([]model.Order, error)
	GetHistoryOrders(id int, p model.Paging) ([]model.Order, error)
	DeleteOrder(id int) error
}
