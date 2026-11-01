package service

import (
	"context"
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
)

type OrderRepository interface {
	CreateOrder(data *model.Order, ctx context.Context) (int, error)
	UpdateOrder(data model.Order, ctx context.Context) error
	GetOrderById(id int, ctx context.Context) (model.Order, error)
	GetOrders(p model.Paging, ctx context.Context) ([]model.Order, error)
	GetHistoryOrders(id int, p model.Paging, ctx context.Context) ([]model.Order, error)
	CancelOrder(id int, ctx context.Context) (restocked []string, err error)
	CancelExpiredOrders(createdBefore time.Time, ctx context.Context) (ids []int, restocked []string, err error)
}

type StockRepository interface {
	ReserveStock(productId string, qnantity int, orderId int, ctx context.Context) (remaining int, err error)
	ReleaseStock(orderId int, ctx context.Context) (restocked []string, err error)
	AddStock(productId string, quantity int, ctx context.Context) (int, error)
	GetStocks(productIds []string, ctx context.Context) ([]model.Stock, error)
}

// StockCache holds the last known quantity of products an order found short, so later
// orders asking for more can be refused before the catalog and the database are called.
type StockCache interface {
	Set(productId string, quantity int, ctx context.Context) error
	Delete(productIds []string, ctx context.Context) error
	// Get returns the cached quantities of productIds; ids that aren't cached are left out.
	Get(productIds []string, ctx context.Context) (map[string]int, error)
}

type SequenecerRepository interface {
	GetOrderId(ctx context.Context) (int, error)
}

type PaymentRepository interface {
	CreateOrderPayment(data model.OrderPayment, ctx context.Context) (int, error)
	UpdateOrderPayment(data model.OrderPayment, ctx context.Context) error
	CompletePayment(eventId string, orderId int, sessionId string, ctx context.Context) error
	FailPayment(eventId string, orderId int, sessionId string, ctx context.Context) error
	GetOrderPayment(orderId int, ctx context.Context) (model.OrderPayment, error)
}
