package model

type OrderPayment struct {
	OrderId           int
	UserId            int
	CheckoutSessionId string
	State             PaymentState
	PaymentIntentId   int
	AmountTotal       int
	Currency          string
}

type PaymentState int

const (
	IsCreated PaymentState = iota
	IsProcessing
	IsCompleted
	IsFailed
)
