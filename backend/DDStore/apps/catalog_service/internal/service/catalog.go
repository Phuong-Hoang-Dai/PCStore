package service

import (
	"context"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type productManager struct {
	repository ProductRepos
}

func NewProductService(repos ProductRepos) ProductService {
	if repos == nil {
		panic("to gen an error message")
	}

	return productManager{
		repository: repos,
	}
}

func (service productManager) CreateProduct(ctx context.Context, data model.Product) (bson.ObjectID, error) {
	if id, err := service.repository.CreateProduct(ctx, data); err != nil {
		return bson.NilObjectID, err
	} else {
		return id, nil
	}
}

func (service productManager) UpdateProduct(ctx context.Context, data model.Product) error {
	return service.repository.UpdateProduct(ctx, data)
}

func (service productManager) GetProducts(ctx context.Context, p *model.Paging) ([]model.Product, error) {
	p.Process()
	return service.repository.GetProducts(ctx, *p)
}

func (service productManager) GetProductsByCate(ctx context.Context, p *model.Paging, cate model.Category) ([]model.Product, error) {
	p.Process()
	return service.repository.GetProductsByCate(ctx, *p, cate)
}

func (service productManager) GetProductById(ctx context.Context, id bson.ObjectID) (model.Product, error) {
	return service.repository.GetProductById(ctx, id)
}

func (service productManager) DeleteProduct(ctx context.Context, id bson.ObjectID) error {
	return service.repository.DeleteProduct(ctx, id)
}