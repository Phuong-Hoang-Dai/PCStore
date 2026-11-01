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
	UpdateProductState(ctx context.Context, id bson.ObjectID, state model.ProductState) error
	RenameCategory(ctx context.Context, cate model.ProductCategory) error
	CountProductsByCate(ctx context.Context, cateId int) (int64, error)
	IsProductExist(ctx context.Context, id bson.ObjectID) (bool, error)
	GetProductById(ctx context.Context, id bson.ObjectID) (model.Product, error)
	GetProductsByIds(ctx context.Context, ids []bson.ObjectID) ([]model.Product, error)
	GetProducts(ctx context.Context, p model.Paging) ([]model.Product, error)
	GetProductsByCate(ctx context.Context, p model.Paging, cateId int) ([]model.Product, error)
	DeleteProduct(ctx context.Context, id bson.ObjectID) error
}

type CateRepos interface {
	CreateCate(ctx context.Context, data model.Category) (int, error)
	UpdateCate(ctx context.Context, data model.Category) error
	GetCateById(ctx context.Context, id int) (model.Category, error)
	GetCates(ctx context.Context) ([]model.Category, error)
	DeleteCate(ctx context.Context, id int) error
}

// ProductEvents tells other services (search) about product changes.
type ProductEvents interface {
	ProductUpserted(ctx context.Context, product model.Product) error
	ProductDeleted(ctx context.Context, id bson.ObjectID) error
}

type ProductService interface {
	CreateProduct(ctx context.Context, data model.Product) (bson.ObjectID, error)
	UpdateProduct(ctx context.Context, data model.Product) error
	UpdateProductState(ctx context.Context, id bson.ObjectID, state model.ProductState) error
	GetProducts(ctx context.Context, p *model.Paging) ([]model.Product, error)
	GetProductsByCate(ctx context.Context, p *model.Paging, cateId int) ([]model.Product, error)
	GetProductById(ctx context.Context, id bson.ObjectID) (model.Product, error)
	GetProductsByIds(ctx context.Context, ids []bson.ObjectID) ([]model.Product, error)
	GetPrices(ctx context.Context, items []OrderItemsDto) ([]OrderItemsDto, error)
	DeleteProduct(ctx context.Context, id bson.ObjectID) error
}

type CateService interface {
	CreateCate(ctx context.Context, data model.Category) (int, error)
	UpdateCate(ctx context.Context, data model.Category) error
	GetCateById(ctx context.Context, id int) (model.Category, error)
	GetCates(ctx context.Context) ([]model.Category, error)
	DeleteCate(ctx context.Context, id int) error
}
