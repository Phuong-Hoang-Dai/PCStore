package service

import (
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
)

type MockRepos struct{}

var productList []model.Product

func (MockRepos) Init() {
	productList = []model.Product{}
}

func (MockRepos) CreateProduct(data model.Product) (int, error) {
	productList = append(productList, data)
	productList[len(productList)-1].Id = len(productList)

	return productList[len(productList)-1].Id, nil
}

func (MockRepos) UpdateProduct(data model.Product) error {
	productList[data.Id-1] = data

	return nil
}

func (MockRepos) UpdateProducts(data []model.Product) error {
	for i := range data {
		productList[data[i].Id-1] = data[i]
	}

	return nil
}

func (MockRepos) GetProductById(id int) (model.Product, error) {
	return productList[id-1], nil
}

func (MockRepos) GetProducts(p model.Paging) (data []model.Product, err error) {
	for i := p.Offset; i < len(productList); i++ {
		if i-p.Offset < p.Limit {
			data = append(data, productList[i])
		}
	}

	return data, nil
}

func (MockRepos) GetProductsByCate(p model.Paging, cate model.Category) (data []model.Product, err error) {
	for i := p.Offset; i < len(productList); i++ {
		if i-p.Offset < p.Limit {
			data = append(data, productList[i])
		}
	}

	return data, nil
}

func (MockRepos) DeleteProduct(id int) error {
	return nil
}
