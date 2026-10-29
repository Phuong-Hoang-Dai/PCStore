package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/repos"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StockHandler struct {
	stockService service.StockService
}

func InitStock(db *pgxpool.Pool, productUrl string) StockHandler {
	return StockHandler{
		stockService: service.StockService{
			Repos:          repos.NewPgStockRepo(db),
			ProductService: productHTTPClient{baseURL: fmt.Sprint(productUrl, "/product")},
		},
	}
}

// GetStocks handles GET /stock?ids=a,b,c; without ids it returns every product's stock.
func (s StockHandler) GetStocks() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var ids []string
		for _, id := range strings.Split(ctx.Query("ids"), ",") {
			if id = strings.TrimSpace(id); id != "" {
				ids = append(ids, id)
			}
		}

		data, err := s.stockService.GetStocks(ids, ctx)
		if err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Stock retrieved successfully",
			"data":    data,
		})
	}
}

// AddStock handles POST /stock/:productId/add with {"quantity": n}.
func (s StockHandler) AddStock() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		productId := ctx.Param("productId")

		var body struct {
			Quantity int `json:"quantity"`
		}
		if err := ctx.ShouldBindJSON(&body); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		total, err := s.stockService.AddStock(productId, body.Quantity, ctx)
		if err != nil {
			responeOrderError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Stock added successfully",
			"data": gin.H{
				"productId": productId,
				"quantity":  total,
			},
		})
	}
}
