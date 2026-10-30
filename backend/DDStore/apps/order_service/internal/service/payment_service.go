package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
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
	state := model.IsCreated
	id, err := p.Repo.CreateOrderPayment(model.OrderPayment{
		OrderId: orderId,
		UserId:  userId,
		State:   &state,
	}, ctx)

	return id, err
}

// GuestCheckout places and checks out an order for someone who isn't logged in. The user
// service returns the account using the guest's email, creating one without a password if
// needed. Checkout happens right away because a guest has no session to pay with later.
func (p PaymentService) GuestCheckout(guest GuestDTO, items []OrderDTO, ctx context.Context) (orderId int, url string, err error) {
	if strings.TrimSpace(guest.Name) == "" || strings.TrimSpace(guest.Email) == "" {
		return 0, "", model.ErrInvalidGuest
	}

	userId, err := p.UserService.GetOrCreateGuest(guest, ctx)
	if err != nil {
		return 0, "", err
	}

	orderId, err = p.OrderService.CreateOrder(userId, items, ctx)
	if err != nil {
		return 0, "", err
	}

	url, err = p.CreateCheckoutSession(orderId, ctx)
	if err != nil {
		// The guest can't retry checkout for this order, so release its stock now instead of
		// leaving it reserved until the unpaid-order sweep.
		if cancelErr := p.OrderService.CancelOrder(orderId, ctx); cancelErr != nil {
			log.Printf("cancel guest order %d after failed checkout: %v", orderId, cancelErr)
		}
		return 0, "", err
	}

	return orderId, url, nil
}

func (p PaymentService) CreateCheckoutSession(orderId int, ctx context.Context) (url string, err error) {
	paymentOrder, err := p.Repo.GetOrderPayment(orderId, ctx)
	if err != nil {
		return "", err
	}

	order, err := p.OrderService.GetOrderById(orderId, ctx)
	if err != nil {
		return "", err
	}
	// Without this a paid order could be checked out (and charged) again, and the
	// UpdateOrderPayment below would move its completed payment back to processing.
	if order.State != model.Processing && order.State != model.ReadyToPay {
		return "", model.ErrOrderNotPaid
	}

	window := time.Now().Unix() / 1800
	// The session must close before OrderService.CancelExpiredOrders cancels the order,
	// or a customer could pay for an order whose stock was already released.
	expiresAt := (window + 2) * 1800
	if payBy := order.CreatedAt.Add(model.OrderPaymentTimeout - model.CheckoutExpiryMargin).Unix(); payBy < expiresAt {
		expiresAt = payBy
	}
	// Stripe rejects sessions that expire less than 30 minutes from now.
	if expiresAt < time.Now().Add(31*time.Minute).Unix() {
		return "", model.ErrOrderPaymentExpired
	}

	stripeUser, err := p.getOrCreateCustomer(ctx, paymentOrder.UserId, paymentOrder.PaymentInfo.StripeUserId)
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
		ExpiresAt: stripe.Int64(expiresAt),
	}
	params.SetIdempotencyKey(fmt.Sprintf("checkout-%d-%d", orderId, window))

	session, err := p.StripeCLient.V1CheckoutSessions.Create(ctx, params)
	if err != nil {
		return "", err
	}

	processing := model.IsProcessing
	err = p.Repo.UpdateOrderPayment(model.OrderPayment{
		OrderId:     orderId,
		State:       &processing,
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

// HandleWebhookEvent updates the payment and order once Stripe reports the checkout result
// in a snapshot event, which carries the checkout session itself.
// Card payments finish with checkout.session.completed; delayed methods (bank transfer, ...)
// complete later with async_payment_succeeded / async_payment_failed.
func (p PaymentService) HandleWebhookEvent(event stripe.Event, ctx context.Context) error {
	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted,
		stripe.EventTypeCheckoutSessionAsyncPaymentSucceeded,
		stripe.EventTypeCheckoutSessionAsyncPaymentFailed,
		stripe.EventTypeCheckoutSessionExpired:
	default:
		return nil
	}

	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		return err
	}
	return p.applyCheckoutResult(event.ID, event.Type, &session, ctx)
}

// HandleEventNotification does the same for a thin event, which only references the checkout
// session, so the session is fetched from Stripe.
func (p PaymentService) HandleEventNotification(notif stripe.EventNotificationContainer, ctx context.Context) error {
	var source interface {
		FetchRelatedObject(ctx context.Context) (*stripe.CheckoutSession, error)
	}
	var eventType stripe.EventType
	switch n := notif.(type) {
	case *stripe.V1CheckoutSessionCompletedEventNotification:
		source, eventType = n, stripe.EventTypeCheckoutSessionCompleted
	case *stripe.V1CheckoutSessionAsyncPaymentSucceededEventNotification:
		source, eventType = n, stripe.EventTypeCheckoutSessionAsyncPaymentSucceeded
	case *stripe.V1CheckoutSessionAsyncPaymentFailedEventNotification:
		source, eventType = n, stripe.EventTypeCheckoutSessionAsyncPaymentFailed
	case *stripe.V1CheckoutSessionExpiredEventNotification:
		source, eventType = n, stripe.EventTypeCheckoutSessionExpired
	default:
		return nil
	}

	session, err := source.FetchRelatedObject(ctx)
	if err != nil {
		return err
	}
	return p.applyCheckoutResult(notif.GetEventNotification().ID, eventType, session, ctx)
}

func (p PaymentService) applyCheckoutResult(eventId string, eventType stripe.EventType, session *stripe.CheckoutSession, ctx context.Context) error {
	orderId, err := strconv.Atoi(session.ClientReferenceID)
	if err != nil {
		return fmt.Errorf("invalid client_reference_id %q: %w", session.ClientReferenceID, err)
	}

	switch eventType {
	case stripe.EventTypeCheckoutSessionCompleted:
		if session.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
			return nil
		}
		return p.Repo.CompletePayment(eventId, orderId, session.ID, ctx)
	case stripe.EventTypeCheckoutSessionAsyncPaymentSucceeded:
		return p.Repo.CompletePayment(eventId, orderId, session.ID, ctx)
	default:
		return p.Repo.FailPayment(eventId, orderId, session.ID, ctx)
	}
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
