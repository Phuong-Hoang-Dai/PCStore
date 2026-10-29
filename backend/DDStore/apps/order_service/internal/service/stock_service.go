package service

import (
	"context"
	"errors"
	"log"

	"github.com/Phuong-Hoang-Dai/DDStore/app/model/catalog"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
)

// StockService owns stock quantities and keeps the catalog's in/out-of-stock flag in step.
// Catalog updates are best effort: a failure is logged, never returned, so an
// unreachable catalog can't block orders.
type StockService struct {
	Repos          StockRepository
	ProductService ProductService
}

func (svc StockService) ReserveStock(orderId int, productId string, quantity int, ctx context.Context) error {
	err := svc.Repos.ReserveStock(productId, quantity, orderId, ctx)

	if errors.Is(err, model.ErrOutOfStock) {
		svc.setState(productId, catalog.OutOfStock, ctx)
	}

	if err != nil {
		return err
	}

	return nil
}

func (svc StockService) ReleaseStock(orderId int, ctx context.Context) error {
	products, err := svc.Repos.ReleaseStock(orderId, ctx)
	svc.MarkAvailable(products, ctx)
	return err
}

// AddStock adds quantity to a product that exists in the catalog and returns the new total.
func (svc StockService) AddStock(productId string, quantity int, ctx context.Context) (int, error) {
	if quantity <= 0 {
		return 0, model.ErrInvalidQuantity
	}

	products, err := svc.ProductService.GetProductsByIds([]string{productId}, ctx)
	if err != nil {
		return 0, err
	}
	if len(products) == 0 {
		return 0, model.ErrProductNotFound
	}

	total, err := svc.Repos.AddStock(productId, quantity, ctx)
	if err != nil {
		return 0, err
	}

	svc.MarkAvailable([]string{productId}, ctx)
	return total, nil
}

func (svc StockService) GetStocks(productIds []string, ctx context.Context) ([]model.Stock, error) {
	return svc.Repos.GetStocks(productIds, ctx)
}

func (svc StockService) MarkAvailable(productIds []string, ctx context.Context) {
	for _, id := range productIds {
		svc.setState(id, catalog.IsAvailable, ctx)
	}
}

func (svc StockService) setState(productId string, state catalog.ProductState, ctx context.Context) {
	if err := svc.ProductService.UpdateProductState(productId, state, ctx); err != nil {
		log.Printf("Failed to set product %s state to %d: %v", productId, state, err)
	}
}
