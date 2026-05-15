package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type productManager struct {
	repository ProductRepos
	rdb        *redis.Client
}

func NewProductService(repos ProductRepos, redis *redis.Client) ProductService {
	if repos == nil {
		panic("to gen an error message")
	}

	return productManager{
		repository: repos,
		rdb:        redis,
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

func (service productManager) GetProducts(ctx context.Context, p *model.Paging) (products []model.Product, err error) {
	cacheKey := fmt.Sprintf("products:%v-%v", p.Limit, p.Offset)

	data, err := service.rdb.Get(ctx, cacheKey).Bytes()
	if err == nil {
		json.Unmarshal(data, &products)
		return products, nil
	}

	p.Process()
	products, err = service.repository.GetProducts(ctx, *p)
	if err != nil {
		return nil, err
	}

	prs, err := json.Marshal(products)
	service.rdb.Set(ctx, cacheKey, prs, 30*time.Second)

	return products, nil
}

func (service productManager) GetProductsByCate(ctx context.Context, p *model.Paging, cate model.Category) (products []model.Product, err error) {
	cacheKey := fmt.Sprintf("productsByCate:%v-%v-%v", p.Limit, p.Offset, cate.Id)

	data, err := service.rdb.Get(ctx, cacheKey).Bytes()
	if err == nil {
		json.Unmarshal(data, &products)
		return products, nil
	}

	p.Process()
	products, err = service.repository.GetProductsByCate(ctx, *p, cate)
	if err != nil {
		return nil, err
	}

	prs, err := json.Marshal(products)
	service.rdb.Set(ctx, cacheKey, prs, 30*time.Second)

	return products, nil
}

func (service productManager) GetProductById(ctx context.Context, id bson.ObjectID) (product model.Product, err error) {
	cacheKey := fmt.Sprintf("product-%v", id)

	data, err := service.rdb.Get(ctx, cacheKey).Bytes()
	if err == nil {
		json.Unmarshal(data, &product)
		return product, nil
	}

	product, err = service.repository.GetProductById(ctx, id)
	if err != nil {
		return product, err
	}

	pr, err := json.Marshal(product)
	service.rdb.Set(ctx, cacheKey, pr, 30*time.Second)

	return product, nil
}

func (service productManager) DeleteProduct(ctx context.Context, id bson.ObjectID) error {
	return service.repository.DeleteProduct(ctx, id)
}
