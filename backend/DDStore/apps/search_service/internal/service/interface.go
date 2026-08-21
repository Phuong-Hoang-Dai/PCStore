package service

import (
	"context"

	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/internal/model"
)

// ProductRepos abstracts the Elasticsearch storage layer so the service
// can be unit-tested with a mock.
type ProductRepos interface {
	EnsureIndex(ctx context.Context) (bool, error)
	Index(ctx context.Context, p model.Product) (string, error)
	BulkIndex(ctx context.Context, products []model.Product) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, req model.SearchRequest) (model.SearchResult, error)
	AutoComplete(ctx context.Context, req model.SearchRequest) (model.AutoCompleteResult, error)
}

// SearchService is the business-facing API consumed by the handlers.
type SearchService interface {
	EnsureIndex(ctx context.Context) (bool, error)
	IndexProduct(ctx context.Context, p model.Product) (string, error)
	IndexProducts(ctx context.Context, products []model.Product) error
	DeleteProduct(ctx context.Context, id string) error
	Search(ctx context.Context, req model.SearchRequest) (model.SearchResult, error)
	AutoComplete(ctx context.Context, req model.SearchRequest) (model.AutoCompleteResult, error)
}
