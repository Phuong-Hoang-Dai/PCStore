package model

type OrderPayment struct {
	OrderId     int
	UserId      int
	State       *PaymentState
	AmountTotal *int
	Currency    *string
	PaymentInfo *Payment
}

type Payment struct {
	StripeUserId      string
	PaymentIntentId   string
	CheckoutSessionId string
}

type PaymentState int

const (
	IsCreated PaymentState = iota
	IsProcessing
	IsCompleted
	IsFailed
)
