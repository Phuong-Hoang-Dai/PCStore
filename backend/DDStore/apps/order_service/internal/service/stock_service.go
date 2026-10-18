package service

import (
	"context"
	"errors"

	"github.com/Phuong-Hoang-Dai/DDStore/app/model/catalog"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
)

type StockService struct {
	Repos          StockRepository
	ProductService ProductService
}

func (svc StockService) ReserveStock(orderId int, productId string, quantity int, ctx context.Context) error {
	err := svc.Repos.ReserveStock(productId, quantity, orderId, ctx)

	if errors.Is(err, model.ErrOutOfStock) {
		svc.ProductService.UpdateProductState(productId, catalog.OutOfStock, ctx)
	}

	return err
}
