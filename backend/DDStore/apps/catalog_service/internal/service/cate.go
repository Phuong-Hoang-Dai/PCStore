package service

import (
	"context"
	"log"
	"strings"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/logger"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type cateManager struct {
	repository CateRepos
	products   ProductRepos
	events     ProductEvents
	cache      productCache
}

func NewCateService(repos CateRepos, products ProductRepos, redis *redis.Client, events ProductEvents) CateService {
	return cateManager{repository: repos, products: products, events: events, cache: productCache{rdb: redis}}
}

func (c cateManager) CreateCate(ctx context.Context, data model.Category) (id int, err error) {
	defer logger.Trace(ctx, "CateService.CreateCate", data)(&id, &err)
	data.Name = strings.TrimSpace(data.Name)
	if data.Name == "" {
		return 0, model.ErrCategoryIsInvalid
	}
	return c.repository.CreateCate(ctx, data)
}

// UpdateCate renames the category, including the copy embedded in its products
// and their search documents.
func (c cateManager) UpdateCate(ctx context.Context, data model.Category) (err error) {
	defer logger.Trace(ctx, "CateService.UpdateCate", data)(&err)
	data.Name = strings.TrimSpace(data.Name)
	if data.Name == "" {
		return model.ErrCategoryIsInvalid
	}
	if err := c.repository.UpdateCate(ctx, data); err != nil {
		return err
	}

	if err := c.products.RenameCategory(ctx, model.ProductCategory{Id: data.Id, Name: data.Name}); err != nil {
		return err
	}

	// A zero limit means no limit.
	products, err := c.products.GetProductsByCate(ctx, model.Paging{}, data.Id)
	if err != nil {
		return err
	}
	ids := make([]bson.ObjectID, len(products))
	for i, p := range products {
		ids[i] = p.Id
	}
	c.cache.invalidateProducts(ctx, ids...)
	for _, p := range products {
		if err := c.events.ProductUpserted(ctx, p); err != nil {
			log.Printf("Failed to publish product %s for indexing: %v", p.Id.Hex(), err)
		}
	}
	return nil
}

func (c cateManager) GetCates(ctx context.Context) (cates []model.Category, err error) {
	defer logger.Trace(ctx, "CateService.GetCates")(&cates, &err)
	return c.repository.GetCates(ctx)
}

func (c cateManager) GetCateById(ctx context.Context, id int) (cate model.Category, err error) {
	defer logger.Trace(ctx, "CateService.GetCateById", id)(&cate, &err)
	return c.repository.GetCateById(ctx, id)
}

// DeleteCate refuses to delete a category that products still point to.
func (c cateManager) DeleteCate(ctx context.Context, id int) (err error) {
	defer logger.Trace(ctx, "CateService.DeleteCate", id)(&err)
	count, err := c.products.CountProductsByCate(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return model.ErrCategoryInUse
	}
	return c.repository.DeleteCate(ctx, id)
}
