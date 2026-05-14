package service

import (
	"context"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProductRepos interface {
	CreateProduct(ctx context.Context, data model.Product) (bson.ObjectID, error)
	UpdateProduct(ctx context.Context, data model.Product) error
	UpdateProducts(ctx context.Context, data []model.Product) error
	GetProductById(ctx context.Context, id bson.ObjectID) (model.Product, error)
	GetProducts(ctx context.Context, p model.Paging) ([]model.Product, error)
	GetProductsByCate(ctx context.Context, p model.Paging, cate model.Category) ([]model.Product, error)
	DeleteProduct(ctx context.Context, id bson.ObjectID) error
}

type CateRepos interface {
	CreateCate(data model.Category) (int, error)
	UpdateCate(data model.Category) error
	GetCateById(id int) (model.Category, error)
	GetCates() ([]model.Category, error)
	DeleteCate(id int) error
}

type ProductService interface {
	CreateProduct(ctx context.Context, data model.Product) (bson.ObjectID, error)
	UpdateProduct(ctx context.Context, data model.Product) error
	GetProducts(ctx context.Context, p *model.Paging) ([]model.Product, error)
	GetProductsByCate(ctx context.Context, p *model.Paging, cate model.Category) ([]model.Product, error)
	GetProductById(ctx context.Context, id bson.ObjectID) (model.Product, error)
	DeleteProduct(ctx context.Context, id bson.ObjectID) error
}
