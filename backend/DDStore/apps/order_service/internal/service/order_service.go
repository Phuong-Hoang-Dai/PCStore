package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/logger"
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
	stockCache       StockCache
}

func NewOrderService(repos OrderRepository, productSvc ProductService, stockSvc StockService, sequencerSvc SequencerService, paymentSvc PaymentCreator, stockCache StockCache) OrderService {
	return OrderService{
		repos:            repos,
		productService:   productSvc,
		stockService:     stockSvc,
		sequencerService: sequencerSvc,
		paymentService:   paymentSvc,
		stockCache:       stockCache,
	}
}

func (svc OrderService) CreateOrder(userId int, items []OrderDTO, ctx context.Context) (orderId int, err error) {
	defer logger.Trace(ctx, "OrderService.CreateOrder", userId, items)(&orderId, &err)
	if len(items) == 0 {
		return 0, model.ErrEmptyOrder
	}
	for _, item := range items {
		// A negative quantity would add stock back and make the total negative.
		if item.Quantity <= 0 {
			return 0, model.ErrInvalidQuantity
		}
	}
	// Refuse stock known to be short before calling the catalog or the database.
	if err := svc.CheckInStock(items, ctx); err != nil {
		return 0, err
	}

	if err := svc.productService.GetPriceProduct(&items, ctx); err != nil {
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
		item := &order.Items[i]
		item.OrderId = order.Id
		remaining, err := svc.stockService.ReserveStock(order.Id, item.ProductId, item.Quantity, ctx)
		if errors.Is(err, model.ErrOutOfStock) {
			// This order took the last units: the reservation is kept.
			svc.cacheStock(item.ProductId, 0, ctx)
			continue
		}
		if err != nil {
			if errors.Is(err, model.ErrNotEnoughStock) {
				svc.cacheStock(item.ProductId, remaining, ctx)
			}
			// Give back what earlier items already reserved.
			svc.releaseStock(order.Id, ctx)
			return 0, err
		}
	}

	id, err := svc.repos.CreateOrder(&order, ctx)
	if err != nil {
		svc.releaseStock(order.Id, ctx)
		return 0, err
	}

	if _, err := svc.paymentService.CreatePaymentInfo(id, userId, ctx); err != nil {
		return 0, err
	}

	return id, nil
}

// CheckInStock fails if an item asks for more than its product's cached quantity. If the
// cache can't be read the order goes on: ReserveStock still checks the database.
func (svc OrderService) CheckInStock(items []OrderDTO, ctx context.Context) (err error) {
	defer logger.Trace(ctx, "OrderService.CheckInStock", items)(&err)
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ProductId
	}
	cached, err := svc.stockCache.Get(ids, ctx)
	if err != nil {
		log.Printf("Failed to read the stock cache: %v", err)
		return nil
	}

	for _, item := range items {
		quantity, ok := cached[item.ProductId]
		if !ok || item.Quantity <= quantity {
			continue
		}
		if quantity == 0 {
			return fmt.Errorf("product %s: %w", item.ProductId, model.ErrOutOfStock)
		}
		return fmt.Errorf("product %s: %w", item.ProductId, model.ErrNotEnoughStock)
	}
	return nil
}

// cacheStock remembers the quantity an order found, so orders asking for more fail early.
func (svc OrderService) cacheStock(productId string, quantity int, ctx context.Context) {
	if err := svc.stockCache.Set(productId, quantity, ctx); err != nil {
		log.Printf("Failed to cache the stock of product %s: %v", productId, err)
	}
}

// uncacheStock drops cached quantities that a restock made too low.
func (svc OrderService) uncacheStock(productIds []string, ctx context.Context) {
	if err := svc.stockCache.Delete(productIds, ctx); err != nil {
		log.Printf("Failed to remove products %v from the stock cache: %v", productIds, err)
	}
}

// releaseStock gives back what an order that couldn't be created had reserved.
func (svc OrderService) releaseStock(orderId int, ctx context.Context) {
	restocked, _ := svc.stockService.ReleaseStock(orderId, ctx)
	svc.uncacheStock(restocked, ctx)
}

// CancelOrder cancels an unpaid order and gives its stock back.
func (svc OrderService) CancelOrder(id int, ctx context.Context) (err error) {
	defer logger.Trace(ctx, "OrderService.CancelOrder", id)(&err)
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
	svc.uncacheStock(restocked, ctx)
	return nil
}

// CancelExpiredOrders cancels the orders left unpaid for OrderPaymentTimeout and
// releases their stock. It returns how many orders were cancelled.
func (svc OrderService) CancelExpiredOrders(ctx context.Context) (cancelled int, err error) {
	defer logger.Trace(ctx, "OrderService.CancelExpiredOrders")(&cancelled, &err)
	ids, restocked, err := svc.repos.CancelExpiredOrders(time.Now().Add(-model.OrderPaymentTimeout), ctx)
	if err != nil {
		return 0, err
	}
	svc.stockService.MarkAvailable(restocked, ctx)
	svc.uncacheStock(restocked, ctx)
	return len(ids), nil
}

// RunExpiredOrderCanceller runs CancelExpiredOrders every interval until ctx is done.
// It is safe to run in several replicas: each order is cancelled by one UPDATE only.
func (svc OrderService) RunExpiredOrderCanceller(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		// Each sweep gets its own transaction id, so its logs can be found together.
		sweepCtx := logger.WithTransactionID(ctx, logger.NewTransactionID())
		if n, err := svc.CancelExpiredOrders(sweepCtx); err != nil {
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
	defer logger.Trace(ctx, "OrderService.GetOrders", p)(&data, &err)
	p.Process()
	if data, err = svc.repos.GetOrders(*p, ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (svc OrderService) GetOrderById(id int, ctx context.Context) (order model.Order, err error) {
	defer logger.Trace(ctx, "OrderService.GetOrderById", id)(&order, &err)
	if data, err := svc.repos.GetOrderById(id, ctx); err != nil {
		return data, err
	} else {
		return data, nil
	}
}

func (svc OrderService) UpdateOrder(id int, state model.OrderState, ctx context.Context) (err error) {
	defer logger.Trace(ctx, "OrderService.UpdateOrder", id, state)(&err)
	data := model.Order{Id: id, State: state}
	if err := svc.repos.UpdateOrder(data, ctx); err != nil {
		return err
	}

	return nil
}

func (svc OrderService) GetHistoryOrders(id int, p *model.Paging, ctx context.Context) (data []model.Order, err error) {
	defer logger.Trace(ctx, "OrderService.GetHistoryOrders", id, p)(&data, &err)
	p.Process()
	if data, err = svc.repos.GetHistoryOrders(id, *p, ctx); err != nil {
		return nil, err
	}
	return data, nil
}
