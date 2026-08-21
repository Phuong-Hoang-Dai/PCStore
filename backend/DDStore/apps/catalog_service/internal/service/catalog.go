package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	isValid, err := validateProduct(ctx, data, service.repository)
	if !isValid {
		return bson.NilObjectID, err
	}

	if id, err := service.repository.CreateProduct(ctx, data); err != nil {
		return bson.NilObjectID, err
	} else {
		return id, nil
	}
}

func (service productManager) UpdateProduct(ctx context.Context, data model.Product) error {
	isValid, err := validateProduct(ctx, data, service.repository)
	if !isValid {
		return err
	}

	return service.repository.UpdateProduct(ctx, data)
}

func (service productManager) GetProducts(ctx context.Context, p *model.Paging) (products []model.Product, err error) {
	cacheKey := fmt.Sprintf("products:%v-%v", p.Limit, p.Offset)

	data, err := service.rdb.Get(ctx, cacheKey).Bytes()
	if err == nil {
		log.Print("Cache hit")
		json.Unmarshal(data, &products)
		return products, nil
	}

	p.Process()
	products, err = service.repository.GetProducts(ctx, *p)
	if err != nil {
		return nil, err
	}

	prs, err := json.Marshal(products)
	if err == nil {
		service.rdb.Set(ctx, cacheKey, prs, 30*time.Second)
	}

	return products, nil
}

func (service productManager) GetProductsByCate(ctx context.Context, p *model.Paging, cate model.Category) (products []model.Product, err error) {
	cacheKey := fmt.Sprintf("productsByCate:%v-%v-%v", p.Limit, p.Offset, cate.Id)

	data, err := service.rdb.Get(ctx, cacheKey).Bytes()
	if err == nil {
		log.Print("Cache hit")
		json.Unmarshal(data, &products)
		return products, nil
	}

	p.Process()
	products, err = service.repository.GetProductsByCate(ctx, *p, cate)
	if err != nil {
		return nil, err
	}

	prs, err := json.Marshal(products)
	if err == nil {
		service.rdb.Set(ctx, cacheKey, prs, 30*time.Second)
	}

	return products, nil
}

func (service productManager) GetProductById(ctx context.Context, id bson.ObjectID) (product model.Product, err error) {
	cacheKey := fmt.Sprintf("product-%v", id)

	data, err := service.rdb.Get(ctx, cacheKey).Bytes()
	if err == nil {
		log.Print("Cache hit")
		json.Unmarshal(data, &product)
		return product, nil
	}

	product, err = service.repository.GetProductById(ctx, id)
	if err != nil {
		return product, err
	}

	pr, err := json.Marshal(product)
	if err == nil {
		service.rdb.Set(ctx, cacheKey, pr, 30*time.Second)
	}

	return product, nil
}

func (service productManager) DeleteProduct(ctx context.Context, id bson.ObjectID) error {
	return service.repository.DeleteProduct(ctx, id)
}

func validateProduct(ctx context.Context, product model.Product, repo ProductRepos) (bool, error) {
	errList := ""
	if product.Type == model.ProductTypeComposite {
		for _, u := range product.Option {
			for _, p := range u {
				isValid, err := repo.IsProductExist(ctx, p.Id)
				if err != nil {
					log.Fatal(err)
					return false, err
				}
				if !isValid {
					errList = fmt.Sprint(errList, p.Id, ", ")
				}
			}
		}
	}
	if errList != "" {
		err := errors.New(fmt.Sprint(model.ErrProductIsInvalid, "(", errList, ")"))
		return false, err
	}
	return true, nil
}
