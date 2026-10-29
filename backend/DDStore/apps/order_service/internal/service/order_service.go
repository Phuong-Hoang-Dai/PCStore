package service

import (
	"context"
	"log"
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
)

// PaymentCreator is the part of PaymentService that OrderService needs.
// An interface avoids a recursive struct, since PaymentService holds an OrderService.
type PaymentCreator interface {
	CreatePaymentInfo(orderId int, userId int, ctx context.Context) (int, error)
}

type OrderService struct {
	repos            OrderRepository
	productService   ProductService
	stockService     StockService
	sequencerService SequencerService
	paymentService   PaymentCreator
}

func NewOrderService(repos OrderRepository, productSvc ProductService, stockSvc StockService, sequencerSvc SequencerService, paymentSvc PaymentCreator) OrderService {
	return OrderService{
		repos:            repos,
		productService:   productSvc,
		stockService:     stockSvc,
		sequencerService: sequencerSvc,
		paymentService:   paymentSvc,
	}
}

func (svc OrderService) CreateOrder(userId int, items []OrderDTO, ctx context.Context) (orderId int, err error) {
	if err := svc.productService.GetPriceProduct(&items); err != nil {
		return 0, err
	}
	itemDAO := make([]model.OrderItem, len(items), cap(items))
	for i := range items {
		MapOrderDTOtoOrderItem(items[i], &itemDAO[i])
	}

	order := model.Order{Items: itemDAO, State: model.Processing, UserId: userId}
	order.CalcTotal()

	order.Id, err = svc.sequencerService.GetOrderId(ctx)
	if err != nil {
		return 0, err
	}

	for i := range order.Items {
		order.Items[i].OrderId = order.Id
		err = svc.stockService.ReserveStock(order.Id, order.Items[i].ProductId, order.Items[i].Quantity, ctx)
		if err != nil {
			// Give back what earlier items already reserved.
			svc.stockService.ReleaseStock(order.Id, ctx)
			return 0, err
		}
	}

	id, err := svc.repos.CreateOrder(&order, ctx)
	if err != nil {
		svc.stockService.ReleaseStock(order.Id, ctx)
		return 0, err
	}

	if _, err := svc.paymentService.CreatePaymentInfo(id, userId, ctx); err != nil {
		return 0, err
	}

	return id, nil
}

// CancelOrder cancels an unpaid order and gives its stock back.
func (svc OrderService) CancelOrder(id int, ctx context.Context) (err error) {
	data, err := svc.repos.GetOrderById(id, ctx)
	if err != nil {
		return err
	}
	if data.State == model.Cancelled {
		return model.ErrOrderIsCanceled
	}

	restocked, err := svc.repos.CancelOrder(id, ctx)
	if err != nil {
		return err
	}
	svc.stockService.MarkAvailable(restocked, ctx)
	return nil
}

// CancelExpiredOrders cancels the orders left unpaid for OrderPaymentTimeout and
// releases their stock. It returns how many orders were cancelled.
func (svc OrderService) CancelExpiredOrders(ctx context.Context) (int, error) {
	ids, restocked, err := svc.repos.CancelExpiredOrders(time.Now().Add(-model.OrderPaymentTimeout), ctx)
	if err != nil {
		return 0, err
	}
	svc.stockService.MarkAvailable(restocked, ctx)
	return len(ids), nil
}

// RunExpiredOrderCanceller runs CancelExpiredOrders every interval until ctx is done.
// It is safe to run in several replicas: each order is cancelled by one UPDATE only.
func (svc OrderService) RunExpiredOrderCanceller(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		if n, err := svc.CancelExpiredOrders(ctx); err != nil {
			log.Printf("Failed to cancel expired orders: %v", err)
		} else if n > 0 {
			log.Printf("Cancelled %d order(s) unpaid after %s", n, model.OrderPaymentTimeout)
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (svc OrderService) GetOrders(p *model.Paging, ctx context.Context) (data []model.Order, err error) {
	p.Process()
	if data, err = svc.repos.GetOrders(*p, ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (svc OrderService) GetOrderById(id int, ctx context.Context) (model.Order, error) {
	if data, err := svc.repos.GetOrderById(id, ctx); err != nil {
		return data, err
	} else {
		return data, nil
	}
}

func (svc OrderService) UpdateOrder(id int, state model.OrderState, ctx context.Context) error {
	data := model.Order{Id: id, State: state}
	if err := svc.repos.UpdateOrder(data, ctx); err != nil {
		return err
	}

	return nil
}

func (svc OrderService) GetHistoryOrders(id int, p *model.Paging, ctx context.Context) (data []model.Order, err error) {
	p.Process()
	if data, err = svc.repos.GetHistoryOrders(id, *p, ctx); err != nil {
		return nil, err
	}
	return data, nil
}
