package service

import (
	"context"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MockRepos struct{}

var productList []model.Product

func (MockRepos) Init() {
	productList = []model.Product{}
}

func (MockRepos) CreateProduct(_ context.Context, data model.Product) (bson.ObjectID, error) {
	data.Id = bson.NewObjectID()
	productList = append(productList, data)
	return data.Id, nil
}

func (MockRepos) UpdateProduct(_ context.Context, data model.Product) error {
	for i := range productList {
		if productList[i].Id == data.Id {
			productList[i] = data
			return nil
		}
	}
	return model.ErrNotFound
}

func (MockRepos) UpdateProducts(_ context.Context, data []model.Product) error {
	for _, d := range data {
		for i := range productList {
			if productList[i].Id == d.Id {
				productList[i] = d
				break
			}
		}
	}
	return nil
}

func (MockRepos) GetProductById(_ context.Context, id bson.ObjectID) (model.Product, error) {
	for _, p := range productList {
		if p.Id == id {
			return p, nil
		}
	}
	return model.Product{}, model.ErrNotFound
}

func (MockRepos) GetProducts(_ context.Context, p model.Paging) (data []model.Product, err error) {
	for i := p.Offset; i < len(productList); i++ {
		if i-p.Offset < p.Limit {
			data = append(data, productList[i])
		}
	}
	return data, nil
}

func (MockRepos) GetProductsByCate(_ context.Context, p model.Paging, cate model.Category) (data []model.Product, err error) {
	for i := p.Offset; i < len(productList); i++ {
		if i-p.Offset < p.Limit {
			data = append(data, productList[i])
		}
	}
	return data, nil
}

func (MockRepos) DeleteProduct(_ context.Context, id bson.ObjectID) error {
	return nil
}