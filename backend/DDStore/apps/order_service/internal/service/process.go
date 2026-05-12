package service

import (
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
)

func CreateOrder(userId int, items []OrderDTO, repos OrderRepository, productService ProductService) (int, error) {
	if err := productService.GetStock(items); err != nil {
		return 0, err
	}
	if err := productService.GetPriceProduct(&items); err != nil {
		productService.RestoreStock(items)
		return 0, err
	}
	itemDAO := make([]model.OrderItem, len(items), cap(items))
	for i := range items {
		MapOrderDTOtoOrderItem(items[i], &itemDAO[i])
	}

	data := model.Order{Items: itemDAO, State: model.Pending, UserId: userId}
	data.CalcTotal()

	id, err := repos.CreateOrder(&data)
	if err != nil {
		productService.RestoreStock(items)
		return 0, err
	}

	return id, nil
}

func CancelOrder(id int, repos OrderRepository, productService ProductService) (err error) {
	data, err := repos.GetOrderById(id)
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

	if err := productService.RestoreStock(itemDTO); err != nil {
		return err
	}

	data.State = model.Cancelled
	if err := repos.UpdateOrder(data); err != nil {
		productService.GetStock(itemDTO)
		return err
	}

	if err := repos.DeleteOrder(id); err != nil {
		return err
	}

	return nil
}

func GetOrders(p *model.Paging, repos OrderRepository) (data []model.Order, err error) {
	p.Process()
	if data, err = repos.GetOrders(*p); err != nil {
		return nil, err
	}
	return data, nil
}

func GetOrderById(id int, repos OrderRepository) (model.Order, error) {
	if data, err := repos.GetOrderById(id); err != nil {
		return data, err
	} else {
		return data, nil
	}
}

func UpdateOrder(id int, state int, repos OrderRepository) error {
	data := model.Order{Id: id, State: state}
	if err := repos.UpdateOrder(data); err != nil {
		return err
	}

	return nil
}

func GetHistoryOrders(id int, p *model.Paging, repos OrderRepository) (data []model.Order, err error) {
	p.Process()
	if data, err = repos.GetHistoryOrders(id, *p); err != nil {
		return nil, err
	}
	return data, nil
}
