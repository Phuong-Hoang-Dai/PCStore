package service

import (
	"context"

	"github.com/Phuong-Hoang-Dai/DDStore/app/model/catalog"
)

type ProductService interface {
	GetProductsByIds(ids []string, ctx context.Context) ([]catalog.Product, error)
	UpdateProductState(productId string, state catalog.ProductState, ctx context.Context) (err error)
	GetPriceProduct(items *[]OrderDTO, ctx context.Context) error
}
