package service

import (
	"context"
	"fmt"

	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/internal/logger"
	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/internal/model"
)

type searchManager struct {
	repository ProductRepos
}

func NewSearchService(repos ProductRepos) SearchService {
	if repos == nil {
		panic("product repository must not be nil")
	}
	return searchManager{repository: repos}
}

// EnsureIndex creates the products index with its mapping if it does not exist.
func (s searchManager) EnsureIndex(ctx context.Context) (created bool, err error) {
	defer logger.Trace(ctx, "SearchService.EnsureIndex")(&created, &err)
	return s.repository.EnsureIndex(ctx)
}

// IndexProduct inserts (or replaces) a single product document and returns its id.
func (s searchManager) IndexProduct(ctx context.Context, p model.Product) (id string, err error) {
	defer logger.Trace(ctx, "SearchService.IndexProduct", p)(&id, &err)
	if err := validateProduct(p); err != nil {
		return "", err
	}
	return s.repository.Index(ctx, p)
}

// IndexProducts bulk-inserts multiple product documents in a single request.
func (s searchManager) IndexProducts(ctx context.Context, products []model.Product) (err error) {
	defer logger.Trace(ctx, "SearchService.IndexProducts", products)(&err)
	if len(products) == 0 {
		return fmt.Errorf("%w: no products to index", model.ErrInvalidInput)
	}
	for _, p := range products {
		if err := validateProduct(p); err != nil {
			return err
		}
	}
	return s.repository.BulkIndex(ctx, products)
}

func (s searchManager) DeleteProduct(ctx context.Context, id string) (err error) {
	defer logger.Trace(ctx, "SearchService.DeleteProduct", id)(&err)
	return s.repository.Delete(ctx, id)
}

func (s searchManager) Search(ctx context.Context, req model.SearchRequest) (result model.SearchResult, err error) {
	defer logger.Trace(ctx, "SearchService.Search", req)(&result, &err)
	return s.repository.Search(ctx, req)
}

func (s searchManager) AutoComplete(ctx context.Context, req model.SearchRequest) (result model.AutoCompleteResult, err error) {
	defer logger.Trace(ctx, "SearchService.AutoComplete", req)(&result, &err)
	return s.repository.AutoComplete(ctx, req)
}

func validateProduct(p model.Product) error {
	return nil
}
