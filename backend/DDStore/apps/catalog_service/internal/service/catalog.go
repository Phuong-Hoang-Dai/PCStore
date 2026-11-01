package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/logger"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type productManager struct {
	repository ProductRepos
	categories CateRepos
	rdb        *redis.Client
	events     ProductEvents
}

func NewProductService(repos ProductRepos, categories CateRepos, redis *redis.Client, events ProductEvents) ProductService {
	if repos == nil || categories == nil || events == nil {
		panic("to gen an error message")
	}

	return productManager{
		repository: repos,
		categories: categories,
		rdb:        redis,
		events:     events,
	}
}

func (service productManager) CreateProduct(ctx context.Context, data model.Product) (id bson.ObjectID, err error) {
	defer logger.Trace(ctx, "ProductService.CreateProduct", data)(&id, &err)
	if err := service.prepareProduct(ctx, &data); err != nil {
		return bson.NilObjectID, err
	}
	if data.Type == "" {
		data.Type = model.ProductTypeStandard
	}
	// A new product has no stock yet; order_service marks it available once stock is added.
	data.State = model.OutOfStock

	id, err = service.repository.CreateProduct(ctx, data)
	if err != nil {
		return bson.NilObjectID, err
	}

	service.invalidateLists(ctx)
	service.publishUpsert(ctx, id)
	return id, nil
}

func (service productManager) UpdateProduct(ctx context.Context, data model.Product) (err error) {
	defer logger.Trace(ctx, "ProductService.UpdateProduct", data)(&err)
	if err := service.prepareProduct(ctx, &data); err != nil {
		return err
	}

	if err := service.repository.UpdateProduct(ctx, data); err != nil {
		return err
	}

	service.invalidateProduct(ctx, data.Id)
	service.publishUpsert(ctx, data.Id)
	return nil
}

func (service productManager) UpdateProductState(ctx context.Context, id bson.ObjectID, state model.ProductState) (err error) {
	defer logger.Trace(ctx, "ProductService.UpdateProductState", id, state)(&err)
	if state != model.IsAvailable && state != model.OutOfStock {
		return model.ErrInvalidState
	}

	if err := service.repository.UpdateProductState(ctx, id, state); err != nil {
		return err
	}

	service.invalidateProduct(ctx, id)
	service.publishUpsert(ctx, id)
	return nil
}

func (service productManager) GetProducts(ctx context.Context, p *model.Paging) (products []model.Product, err error) {
	defer logger.Trace(ctx, "ProductService.GetProducts", p)(&products, &err)
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

func (service productManager) GetProductsByCate(ctx context.Context, p *model.Paging, cateId int) (products []model.Product, err error) {
	defer logger.Trace(ctx, "ProductService.GetProductsByCate", p, cateId)(&products, &err)
	cacheKey := fmt.Sprintf("productsByCate:%v-%v-%v", p.Limit, p.Offset, cateId)

	data, err := service.rdb.Get(ctx, cacheKey).Bytes()
	if err == nil {
		log.Print("Cache hit")
		json.Unmarshal(data, &products)
		return products, nil
	}

	p.Process()
	products, err = service.repository.GetProductsByCate(ctx, *p, cateId)
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
	defer logger.Trace(ctx, "ProductService.GetProductById", id)(&product, &err)
	cacheKey := productCacheKey(id)

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

func (service productManager) DeleteProduct(ctx context.Context, id bson.ObjectID) (err error) {
	defer logger.Trace(ctx, "ProductService.DeleteProduct", id)(&err)
	if err := service.repository.DeleteProduct(ctx, id); err != nil {
		return err
	}

	service.invalidateProduct(ctx, id)
	if err := service.events.ProductDeleted(ctx, id); err != nil {
		log.Printf("Failed to publish deletion of product %s: %v", id.Hex(), err)
	}
	return nil
}

func (service productManager) GetProductsByIds(ctx context.Context, ids []bson.ObjectID) (products []model.Product, err error) {
	defer logger.Trace(ctx, "ProductService.GetProductsByIds", ids)(&products, &err)
	return service.repository.GetProductsByIds(ctx, ids)
}

func (service productManager) GetPrices(ctx context.Context, items []OrderItemsDto) (priced []OrderItemsDto, err error) {
	defer logger.Trace(ctx, "ProductService.GetPrices", items)(&priced, &err)
	ids := make([]bson.ObjectID, len(items))
	for i := range items {
		id, err := bson.ObjectIDFromHex(items[i].ProductId)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid product id %s", model.ErrNotFound, items[i].ProductId)
		}
		ids[i] = id
	}

	products, err := service.repository.GetProductsByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	prices := make(map[bson.ObjectID]float64, len(products))
	for _, p := range products {
		prices[p.Id] = p.Price
	}

	for i := range items {
		price, ok := prices[ids[i]]
		if !ok {
			return nil, fmt.Errorf("%w: product %s", model.ErrNotFound, items[i].ProductId)
		}
		items[i].Price = price
	}

	return items, nil
}

// prepareProduct validates the product and copies the current category name into it.
func (service productManager) prepareProduct(ctx context.Context, data *model.Product) error {
	if isValid, err := validateProduct(ctx, *data, service.repository); !isValid {
		return err
	}

	cate, err := service.categories.GetCateById(ctx, data.Cate.Id)
	if errors.Is(err, model.ErrNotFound) {
		return fmt.Errorf("%w: category %d doesn't exist", model.ErrProductIsInvalid, data.Cate.Id)
	}
	if err != nil {
		return err
	}
	data.Cate = model.ProductCategory{Id: cate.Id, Name: cate.Name}
	return nil
}

// publishUpsert sends the stored product (with its id and timestamps) to the search index.
// A failed publish is logged, not returned: the product change itself already succeeded.
func (service productManager) publishUpsert(ctx context.Context, id bson.ObjectID) {
	product, err := service.repository.GetProductById(ctx, id)
	if err == nil {
		err = service.events.ProductUpserted(ctx, product)
	}
	if err != nil {
		log.Printf("Failed to publish product %s for indexing: %v", id.Hex(), err)
	}
}

func productCacheKey(id bson.ObjectID) string {
	return fmt.Sprintf("product-%v", id)
}

// productCache drops cached product reads after a write, so changes show up
// immediately instead of after the 30s TTL. Shared by product and category writes.
type productCache struct {
	rdb *redis.Client
}

func (c productCache) invalidateProducts(ctx context.Context, ids ...bson.ObjectID) {
	for _, id := range ids {
		c.rdb.Del(ctx, productCacheKey(id))
	}
	c.invalidateLists(ctx)
}

func (c productCache) invalidateLists(ctx context.Context) {
	for _, pattern := range []string{"products:*", "productsByCate:*"} {
		iter := c.rdb.Scan(ctx, 0, pattern, 100).Iterator()
		for iter.Next(ctx) {
			c.rdb.Del(ctx, iter.Val())
		}
		if err := iter.Err(); err != nil {
			log.Printf("Failed to invalidate %s cache: %v", pattern, err)
		}
	}
}

func (service productManager) invalidateProduct(ctx context.Context, id bson.ObjectID) {
	productCache{rdb: service.rdb}.invalidateProducts(ctx, id)
}

func (service productManager) invalidateLists(ctx context.Context) {
	productCache{rdb: service.rdb}.invalidateLists(ctx)
}

func validateProduct(ctx context.Context, product model.Product, repo ProductRepos) (bool, error) {
	if strings.TrimSpace(product.Name) == "" {
		return false, fmt.Errorf("%w: name is required", model.ErrProductIsInvalid)
	}
	if product.Price < 0 {
		return false, fmt.Errorf("%w: price must not be negative", model.ErrProductIsInvalid)
	}

	errList := ""
	if product.Type == model.ProductTypeComposite {
		for _, u := range product.Option {
			for _, p := range u {
				isValid, err := repo.IsProductExist(ctx, p.Id)
				if err != nil {
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
