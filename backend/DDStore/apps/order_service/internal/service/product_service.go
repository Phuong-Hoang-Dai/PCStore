package service

import (
	"context"

	"github.com/Phuong-Hoang-Dai/DDStore/app/model/catalog"
)

type ProductService interface {
	GetPriceProduct(items *[]OrderDTO) error
	UpdateProductState(productId string, state catalog.ProductState, ctx context.Context) (err error)
}
