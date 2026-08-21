package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/internal/repos"
	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/internal/service"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	service service.SearchService
}

func NewSearchHandler(es *elasticsearch.TypedClient) SearchHandler {
	repo := repos.NewESProductRepo(es, configs.Cfg.ESIndex)
	return SearchHandler{service: service.NewSearchService(repo)}
}

// IndexProduct handles POST /search/products — insert a single product.
func (h SearchHandler) IndexProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.Product
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, BuildResponse("fail to index product", err.Error(), nil))
			return
		}

		id, err := h.service.IndexProduct(c.Request.Context(), data)
		if err != nil {
			httpStatusCode, errMsg := responseServiceError(err)
			c.JSON(httpStatusCode, BuildResponse("fail to index product", errMsg, nil))

			return
		}

		c.JSON(http.StatusOK, BuildResponse("Product indexed successfully", "", gin.H{"id": id}))
	}
}

// IndexProducts handles POST /search/products/bulk — insert many products.
func (h SearchHandler) IndexProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data []model.Product
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, BuildResponse("fail to index products", err.Error(), nil))
			return
		}

		if err := h.service.IndexProducts(c.Request.Context(), data); err != nil {
			httpStatusCode, errMsg := responseServiceError(err)
			c.JSON(httpStatusCode, BuildResponse("fail to index products", errMsg, nil))

			return
		}

		c.JSON(http.StatusOK, BuildResponse("Products indexed successfully", "", gin.H{"count": len(data)}))
	}
}

// DeleteProduct handles DELETE /search/products/:id.
func (h SearchHandler) DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
			httpStatusCode, errMsg := responseServiceError(err)
			c.JSON(httpStatusCode, BuildResponse("fail to delete product", errMsg, nil))

			return
		}

		c.JSON(http.StatusOK, BuildResponse("Product deleted successfully", "", gin.H{"id": id}))
	}
}

// Search handles GET /search/products — full-text search with filters + paging.
//
// Query params: q, brand, type, cate_id, min_price, max_price, limit, offset.
func (h SearchHandler) Search() gin.HandlerFunc {
	return func(c *gin.Context) {

		var req model.SearchRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, BuildResponse("fail to uncode", err.Error(), nil))

			return
		}
		result, err := h.service.Search(c.Request.Context(), req)
		if err != nil {
			httpStatusCode, errMsg := responseServiceError(err)
			c.JSON(httpStatusCode, BuildResponse("fail to search", errMsg, nil))

			return
		}

		response := BuildResponse("Search completed successfully", "", result.Products)
		response["pagination"] = gin.H{
			"offset": req.State.Current + 1,
			"limit":  req.State.ResultsPerPage,
			"total":  result.Total,
		}
		response["facets"] = gin.H{
			"value": result.Facets.Value,
			"range": result.Facets.Range,
		}

		c.JSON(http.StatusOK, response)
	}
}

func (h SearchHandler) AutoComplete() gin.HandlerFunc {
	return func(c *gin.Context) {

		var req model.SearchRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, BuildResponse("fail to uncode", err.Error(), nil))

			return
		}
		result, err := h.service.AutoComplete(c.Request.Context(), req)
		if err != nil {
			httpStatusCode, errMsg := responseServiceError(err)
			c.JSON(httpStatusCode, BuildResponse("fail to search", errMsg, nil))

			return
		}

		response := BuildResponse("Search completed successfully", "", result.Products)
		response["pagination"] = gin.H{
			"offset": req.State.Current + 1,
			"limit":  req.State.ResultsPerPage,
			"total":  result.Total,
		}
		c.JSON(http.StatusOK, response)
	}
}

// responseServiceError maps domain errors to HTTP status codes.
func responseServiceError(err error) (httpStatusCode int, errMsg string) {
	switch {
	case errors.Is(err, model.ErrNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, model.ErrInvalidInput):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, err.Error()
	}
}
