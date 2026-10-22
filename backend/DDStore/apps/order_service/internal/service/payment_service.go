package service

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/app/model/catalog"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
	"github.com/stripe/stripe-go/v87"
)

type PaymentService struct {
	Repo           PaymentRepository
	UserService    UserService
	StripeCLient   *stripe.Client
	OrderService   OrderService
	ProductService ProductService
	SuccessURL     string
	CancelURL      string
}

func (p PaymentService) CreatePaymentInfo(orderId int, userId int, ctx context.Context) (status int, err error) {
	id, err := p.Repo.CreateOrderPayment(model.OrderPayment{
		OrderId: orderId,
		UserId:  userId,
	}, ctx)

	return id, err
}

func (p PaymentService) CreateCheckoutSession(orderId int, ctx context.Context) (url string, err error) {
	paymentOrder, err := p.Repo.GetOrderPayment(orderId, ctx)
	if err != nil {
		return "", err
	}
	stripeUser, err := p.getOrCreateCustomer(ctx, paymentOrder.UserId, paymentOrder.PaymentInfo.StripeUserId)
	if err != nil {
		return "", err
	}

	order, err := p.OrderService.GetOrderById(orderId, ctx)
	if err != nil {
		return "", err
	}
	products, err := p.getProducts(ctx, order.Items)
	if err != nil {
		return "", err
	}

	lineItems := []*stripe.CheckoutSessionCreateLineItemParams{}
	for _, item := range order.Items {
		i := slices.IndexFunc(products, func(p catalog.Product) bool {
			return p.Id.Hex() == item.ProductId
		})
		if i == -1 {
			return "", fmt.Errorf("product %s not found", item.ProductId)
		}
		product := products[i]

		productData := &stripe.CheckoutSessionCreateLineItemPriceDataProductDataParams{
			Name:        stripe.String(product.Name),
			Description: stripe.String(product.Desc),
		}
		if len(product.Images) > 0 {
			productData.Images = []*string{stripe.String(product.Images[0].Url)}
		}

		lineItems = append(lineItems, &stripe.CheckoutSessionCreateLineItemParams{
			PriceData: &stripe.CheckoutSessionCreateLineItemPriceDataParams{
				Currency:    stripe.String(stripe.CurrencyVND),
				ProductData: productData,
				UnitAmount:  stripe.Int64(int64(item.Price)),
			},
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	window := time.Now().Unix() / 1800
	orderRef := strconv.Itoa(orderId)

	params := &stripe.CheckoutSessionCreateParams{
		Customer:          &stripeUser.ID,
		LineItems:         lineItems,
		Mode:              stripe.String(stripe.CheckoutSessionModePayment),
		UIMode:            stripe.String(stripe.CheckoutSessionUIModeHostedPage),
		SuccessURL:        stripe.String(p.SuccessURL),
		CancelURL:         stripe.String(p.CancelURL),
		ClientReferenceID: stripe.String(orderRef),
		Metadata:          map[string]string{"order_id": orderRef},
		PaymentIntentData: &stripe.CheckoutSessionCreatePaymentIntentDataParams{
			Metadata: map[string]string{"order_id": orderRef},
		},
		ExpiresAt: stripe.Int64((window + 2) * 1800),
	}
	params.SetIdempotencyKey(fmt.Sprintf("checkout-%d-%d", orderId, window))

	session, err := p.StripeCLient.V1CheckoutSessions.Create(ctx, params)
	if err != nil {
		return "", err
	}

	err = p.Repo.UpdateOrderPayment(model.OrderPayment{
		OrderId:     orderId,
		AmountTotal: paymentOrder.AmountTotal,
		Currency:    stripe.String(stripe.CurrencyVND),
		PaymentInfo: &model.Payment{
			StripeUserId:      stripeUser.ID,
			CheckoutSessionId: session.ID,
		},
	}, ctx)
	if err != nil {
		return "", err
	}

	return session.URL, nil
}

func (p PaymentService) getOrCreateCustomer(ctx context.Context, userId int, stripeUserId string) (stripeUser *stripe.Customer, err error) {
	user, err := p.UserService.GetUserById(userId, ctx)
	if err != nil {
		return stripeUser, err
	}

	if stripeUserId != "" {
		return p.StripeCLient.V1Customers.Retrieve(ctx, stripeUserId, &stripe.CustomerRetrieveParams{})
	}

	return p.StripeCLient.V1Customers.Create(ctx, &stripe.CustomerCreateParams{
		Name:  &user.Name,
		Email: &user.Email,
	})

}

func (p PaymentService) getProducts(ctx context.Context, items []model.OrderItem) ([]catalog.Product, error) {
	var productIds []string
	for _, v := range items {
		productIds = append(productIds, v.ProductId)
	}
	return p.ProductService.GetProductsByIds(productIds, ctx)
}
