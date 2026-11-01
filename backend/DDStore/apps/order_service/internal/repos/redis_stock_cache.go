package repos

import (
	"context"
	"strconv"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/service"
	"github.com/redis/go-redis/v9"
)

type redisStockCache struct {
	rdb *redis.Client
}

func NewRedisStockCache(rdb *redis.Client) service.StockCache {
	return redisStockCache{rdb: rdb}
}

// stockCacheKey must not match the catalog's product-*, products:* or productsByCate:* keys.
func stockCacheKey(productId string) string {
	return "stock:" + productId
}

func stockCacheKeys(productIds []string) []string {
	keys := make([]string, len(productIds))
	for i, id := range productIds {
		keys[i] = stockCacheKey(id)
	}
	return keys
}

func (c redisStockCache) Set(productId string, quantity int, ctx context.Context) error {
	return c.rdb.Set(ctx, stockCacheKey(productId), quantity, model.StockCacheTTL).Err()
}

func (c redisStockCache) Delete(productIds []string, ctx context.Context) error {
	if len(productIds) == 0 {
		return nil
	}
	return c.rdb.Del(ctx, stockCacheKeys(productIds)...).Err()
}

func (c redisStockCache) Get(productIds []string, ctx context.Context) (map[string]int, error) {
	if len(productIds) == 0 {
		return nil, nil
	}
	values, err := c.rdb.MGet(ctx, stockCacheKeys(productIds)...).Result()
	if err != nil {
		return nil, err
	}

	quantities := make(map[string]int)
	for i, v := range values {
		// MGET gives nil for keys that aren't set.
		s, ok := v.(string)
		if !ok {
			continue
		}
		if q, err := strconv.Atoi(s); err == nil {
			quantities[productIds[i]] = q
		}
	}
	return quantities, nil
}
