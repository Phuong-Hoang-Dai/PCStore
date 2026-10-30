package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/repos"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stripe/stripe-go/v87"
	"github.com/stripe/stripe-go/v87/webhook"
)

type PaymentHandler struct {
	paymentService service.PaymentService
	stripeClient   *stripe.Client
	webhookSecret  string
}

func InitPayment(db *pgxpool.Pool, userUrl string, productUrl string) PaymentHandler {
	stripeClient := stripe.NewClient(configs.Cfg.StripeSecretKey)
	p := PaymentHandler{
		paymentService: service.PaymentService{
			Repo:           repos.NewPgPaymentRepo(db),
			UserService:    userHTTPClient{baseURL: fmt.Sprint(userUrl, "/user")},
			ProductService: productHTTPClient{baseURL: fmt.Sprint(productUrl, "/product")},
			StripeCLient:   stripeClient,
			OrderService: service.NewOrderService(
				repos.NewPgOrderRepo(db),
				productHTTPClient{baseURL: fmt.Sprint(productUrl, "/product")},
				service.StockService{Repos: repos.NewPgStockRepo(db), ProductService: productHTTPClient{baseURL: fmt.Sprint(productUrl, "/product")}},
				service.SequencerService{Repos: repos.NewPgSequencerkRepo(db)},
				service.PaymentService{Repo: repos.NewPgPaymentRepo(db)},
			),
			SuccessURL: configs.Cfg.StripeSuccessURL,
			CancelURL:  configs.Cfg.StripeCancelURL,
		},
		stripeClient:  stripeClient,
		webhookSecret: configs.Cfg.StripeWebhookKey,
	}

	return p
}

func (p PaymentHandler) StripeWebhook() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// The signature is computed over the raw body, so it must be read before any binding.
		payload, err := ctx.GetRawData()
		if err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		signature := ctx.GetHeader("Stripe-Signature")

		// Event destinations with a thin payload send "v2.core.event" notifications; snapshot
		// destinations and `stripe listen` send full "event" objects. Both are verified.
		var envelope struct {
			Object string `json:"object"`
		}
		if err := json.Unmarshal(payload, &envelope); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		if envelope.Object == "v2.core.event" {
			notif, parseErr := p.stripeClient.ParseEventNotification(payload, signature, p.webhookSecret)
			if parseErr != nil {
				responeError(http.StatusBadRequest, parseErr, ctx)
				return
			}
			err = p.paymentService.HandleEventNotification(notif, ctx)
		} else {
			event, parseErr := webhook.ConstructEvent(payload, signature, p.webhookSecret)
			if parseErr != nil {
				responeError(http.StatusBadRequest, parseErr, ctx)
				return
			}
			err = p.paymentService.HandleWebhookEvent(event, ctx)
		}

		// A non-2xx response makes Stripe retry the event.
		if err != nil {
			responeError(http.StatusInternalServerError, err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

// GuestCheckout places an order for a customer who isn't logged in and returns the order id
// and its Stripe checkout URL.
func (p PaymentHandler) GuestCheckout() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		bodyData := struct {
			Name  string             `json:"name"`
			Email string             `json:"email"`
			Data  []service.OrderDTO `json:"data"`
		}{}

		if err := ctx.ShouldBind(&bodyData); err != nil {
			responeError(http.StatusBadRequest, err, ctx)
			return
		}

		guest := service.GuestDTO{Name: bodyData.Name, Email: bodyData.Email}
		id, url, err := p.paymentService.GuestCheckout(guest, bodyData.Data, ctx)
		if err != nil {
			responeOrderError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Create Order succesfully",
			"data": gin.H{
				"id":  id,
				"url": url,
			},
		})
	}
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
		if err != nil {
			responeOrderError(err, ctx)
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
