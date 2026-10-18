package service

import (
	"context"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
)

type OrderService struct {
	repos            OrderRepository
	productService   ProductService
	stockService     StockService
	sequencerService SequencerService
}

func NewOrderService(repos OrderRepository, productSvc ProductService, stockSvc StockService, sequencerSvc SequencerService) OrderService {
	return OrderService{
		repos:            repos,
		productService:   productSvc,
		stockService:     stockSvc,
		sequencerService: sequencerSvc,
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

	for _, v := range itemDAO {
		err = svc.stockService.ReserveStock(v.OrderId, v.ProductId, v.Quantity, ctx)
		if err != nil {
			return 0, err
		}
	}

	id, err := svc.repos.CreateOrder(&order, ctx)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (svc OrderService) CancelOrder(id int, ctx context.Context) (err error) {
	data, err := svc.repos.GetOrderById(id, ctx)
	if err != nil {
		return err
	}
	if data.State == model.Cancelled {
		return model.ErrOrderIsCanceled
	}
	itemDTO := make([]OrderDTO, len(data.Items), cap(data.Items))
	for i := range data.Items {
		MapOrderItemtoOrderDTO(data.Items[i], &itemDTO[i])
	}
	//Restore stock

	if err := svc.repos.DeleteOrder(id, ctx); err != nil {
		return err
	}

	return nil
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
