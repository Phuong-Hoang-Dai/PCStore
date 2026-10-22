package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/repos"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stripe/stripe-go/v87"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func InitPayment(db *pgxpool.Pool, userUrl string, productUrl string) PaymentHandler {
	p := PaymentHandler{
		paymentService: service.PaymentService{
			Repo:         repos.NewPgPaymentRepo(db),
			UserService:  userHTTPClient{baseURL: fmt.Sprint(userUrl, "/user")},
			StripeCLient: stripe.NewClient(configs.Cfg.StripeSecretKey),
			OrderService: service.NewOrderService(
				repos.NewPgOrderRepo(db),
				productHTTPClient{baseURL: fmt.Sprint(productUrl, "/product")},
				service.StockService{Repos: repos.NewPgStockRepo(db), ProductService: productHTTPClient{baseURL: fmt.Sprint(productUrl, "/product")}},
				service.SequencerService{Repos: repos.NewPgSequencerkRepo(db)},
			),
			SuccessURL: configs.Cfg.StripeSuccessURL,
			CancelURL:  configs.Cfg.StripeCancelURL,
		},
	}

	return p
}

func (p PaymentHandler) CreatePaymentInfo() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		orderId, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		bodyData := struct {
			UserId int `json:"userId"`
		}{}

		if err := ctx.ShouldBind(&bodyData); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}
		if bodyData.UserId == 0 {
			responeError(http.StatusBadRequest, model.ErrCannotGetUserId, ctx)
			return
		}

		id, err := p.paymentService.CreatePaymentInfo(orderId, bodyData.UserId, ctx)
		if err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Create Payment succesfully",
			"data": gin.H{
				"id": id,
			},
		})
	}
}

func (p PaymentHandler) CreateCheckoutSession() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		orderId, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		url, err := p.paymentService.CreateCheckoutSession(orderId, ctx)
		if errors.Is(err, pgx.ErrNoRows) {
			responeError(http.StatusNotFound, err, ctx)
			return
		}
		if err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Create Checkout Session succesfully",
			"data": gin.H{
				"url": url,
			},
		})
	}
}
