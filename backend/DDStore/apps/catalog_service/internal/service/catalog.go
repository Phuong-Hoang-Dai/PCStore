package service

import (
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
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

func (service productManager) CreateProduct(data model.Product) (int, error) {
	if id, err := service.repository.CreateProduct(data); err != nil {
		return 0, err
	} else {
		return id, err
	}
}

func (service productManager) UpdateProduct(data model.Product) error {
	if err := service.repository.UpdateProduct(data); err != nil {
		return err
	} else {
		return nil
	}
}

func (service productManager) GetProducts(p *model.Paging) (data []model.Product, err error) {
	p.Process()
	if data, err = service.repository.GetProducts(*p); err != nil {
		return nil, err
	}
	return data, nil
}

func (service productManager) GetProductsByCate(p *model.Paging, cate model.Category) (data []model.Product, err error) {
	p.Process()
	if data, err = service.repository.GetProductsByCate(*p, cate); err != nil {
		return nil, err
	}
	return data, nil
}

func (service productManager) GetProductById(id int) (data model.Product, err error) {
	if data, err := service.repository.GetProductById(id); err != nil {
		return data, err
	} else {
		return data, nil
	}
}

func (service productManager) DeleteProduct(id int) error {
	if err := service.repository.DeleteProduct(id); err != nil {
		return err
	} else {
		return nil
	}
}
