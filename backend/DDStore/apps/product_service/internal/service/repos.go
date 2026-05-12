package service

import (
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
)

type ProductRepos interface {
	CreateProduct(data model.Product) (int, error)
	UpdateProduct(data model.Product) error
	UpdateProducts(data []model.Product) error
	GetProductById(id int) (model.Product, error)
	GetProducts(p model.Paging) ([]model.Product, error)
	GetProductsByCate(p model.Paging, cate model.Category) ([]model.Product, error)
	DeleteProduct(id int) error
}

type CateRepos interface {
	CreateCate(data model.Category) (int, error)
	UpdateCate(data model.Category) error
	GetCateById(id int) (model.Category, error)
	GetCates() ([]model.Category, error)
	DeleteCate(id int) error
}
