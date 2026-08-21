package repos

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/internal/service"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/some"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
)

type esProductRepo struct {
	es    *elasticsearch.TypedClient
	index string
}

// NewESProductRepo builds an Elasticsearch-backed product repository.
func NewESProductRepo(es *elasticsearch.TypedClient, index string) service.ProductRepos {
	if es == nil {
		panic("elasticsearch client must not be nil")
	}
	if index == "" {
		panic("elasticsearch index must not be empty")
	}
	return &esProductRepo{es: es, index: index}
}

// EnsureIndex creates the index with productMapping if it does not exist yet.
func (r *esProductRepo) EnsureIndex(ctx context.Context) (bool, error) {
	exists, err := r.es.Indices.Exists(r.index).Do(ctx)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Index inserts or replaces a single product. When p.Id is empty Elasticsearch
// generates one; the effective document id is returned either way.
func (r *esProductRepo) Index(ctx context.Context, p model.Product) (string, error) {
	stampTimes(&p)
	res, err := r.es.Index(r.index).Document(p).Id(p.Id).Refresh(refresh.True).Do(ctx)
	if err != nil {
		log.Println(fmt.Errorf("Index a product: %s", err))
		return "", err
	}
	return res.Id_, nil
}

// BulkIndex indexes many products in a single _bulk request.
func (r *esProductRepo) BulkIndex(ctx context.Context, products []model.Product) error {
	if len(products) == 0 {
		return nil
	}
	bulk := r.es.Bulk()
	for _, v := range products {
		bulk.IndexOp(types.IndexOperation{Index_: &r.index, Id_: &v.Id}, v)
	}
	res, err := bulk.Refresh(refresh.True).Do(ctx)
	if err != nil {
		log.Println(fmt.Errorf("bulk index products: %s", err))
	}
	if res.Errors {
		log.Println(fmt.Errorf("bulk index products: %v", res.Errors))
	}
	return err
}

func (r *esProductRepo) Delete(ctx context.Context, id string) error {
	res, err := r.es.Delete(r.index, id).Refresh(refresh.True).Do(ctx)
	if err != nil {
		log.Println(fmt.Errorf("delete product: %s", res.Result.String()).Error())
	}
	return nil
}

// Search runs a bool query built from the request and returns the matching page.
func (r *esProductRepo) Search(ctx context.Context, req model.SearchRequest) (model.SearchResult, error) {
	limit, offset := req.State.ResultsPerPage, (req.State.Current-1)*req.State.ResultsPerPage
	res, err := r.es.Search().Index(r.index).Request(&search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Must: []types.Query{
					{
						MultiMatch: &types.MultiMatchQuery{
							Query:     req.State.SearchTerm,
							Fields:    []string{"name^3", "brand^2", "description"},
							Fuzziness: some.Int(1),
						},
					},
				},
				Filter: buildFilter(req.State.Filter),
			},
		},
		Size: some.Int(limit),
		From: some.Int(offset),
		Sort: buildSort(req.State.Sort),
		Aggregations: map[string]types.Aggregations{
			"Tern": {
				Terms: &types.TermsAggregation{
					Field: some.String("brand.keyword"),
				},
			},
		},
	}).TrackTotalHits(true).Do(ctx)

	if err != nil {
		return model.SearchResult{}, err
	}

	result := model.SearchResult{
		Total:    res.Hits.Total.Value,
		Products: make([]model.Product, 0, len(res.Hits.Hits)),
	}

	for _, hit := range res.Hits.Hits {
		var p model.Product
		if err := json.Unmarshal(hit.Source_, &p); err != nil {
			return model.SearchResult{}, err
		}
		result.Products = append(result.Products, p)
	}

	result.Facets = buildFacets(res.Aggregations)

	return result, nil
}

func buildFacets(aggs map[string]types.Aggregate) (facets model.FacetResponse) {

	for _, agg := range aggs {
		terms, ok := agg.(*types.StringTermsAggregate)
		if !ok {
			continue
		}

		buckets, ok := terms.Buckets.([]types.StringTermsBucket)
		if !ok {
			continue
		}

		values := make([]model.FacetValueRes, 0, len(buckets))
		for _, b := range buckets {
			key, ok := b.Key.(string)
			if !ok {
				continue
			}
			values = append(values, model.FacetValueRes{
				Field: key,
				Count: int(b.DocCount),
			})
		}
		facets.Value = values
	}

	return facets
}

func (r *esProductRepo) AutoComplete(ctx context.Context, req model.SearchRequest) (model.AutoCompleteResult, error) {
	res, err := r.es.Search().Index(r.index).Request(&search.Request{
		Query: &types.Query{
			MatchBoolPrefix: map[string]types.MatchBoolPrefixQuery{
				"name": {
					Query:     req.State.SearchTerm,
					Fuzziness: 2,
				},
			},
		},
	}).TrackTotalHits(true).Do(ctx)
	if err != nil {
		return model.AutoCompleteResult{}, err
	}

	result := model.AutoCompleteResult{
		Total:    res.Hits.Total.Value,
		Products: make([]model.Product, 0, len(res.Hits.Hits)),
	}

	for _, hit := range res.Hits.Hits {
		var p model.Product
		if err := json.Unmarshal(hit.Source_, &p); err != nil {
			return model.AutoCompleteResult{}, err
		}
		result.Products = append(result.Products, p)
	}

	return result, nil
}

func stampTimes(p *model.Product) {
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	p.UpdatedAt = now
}

func isAsc(direction int) *sortorder.SortOrder {
	if direction == 1 {
		return &sortorder.Asc
	} else {
		return &sortorder.Desc
	}
}

func buildFilter(filterData []model.FilterValue) (filterList []types.Query) {
	for _, v := range filterData {
		for _, fv := range v.FilterValue {
			filterList = append(filterList, types.Query{
				Term: map[string]types.TermQuery{
					v.Name: {Value: fv},
				},
			})
		}
	}
	return filterList
}

func buildSort(sortData model.SortValue) (sort []types.SortCombinations) {
	if sortData.Sort != "" {
		sort = append(sort, types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				sortData.Sort: {Order: isAsc(sortData.SortDirection)},
			},
		})
	}
	return sort
}
