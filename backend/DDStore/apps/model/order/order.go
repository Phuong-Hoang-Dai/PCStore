package order

import (
	"time"
)

type Order struct {
	Id        int         `json:"orderId"`
	State     OrderState  `json:"state"`
	Items     []OrderItem `json:"items"`
	Total     float64     `json:"total"`
	UserId    int         `json:"userId"`
	CreatedAt time.Time   `json:"create_at"`
	UpdatedAt time.Time   `json:"update_at"`
	DeletedAt *time.Time  `json:"delete_at,omitempty"`
}

type OrderItem struct {
	ProductId string     `json:"productId"`
	Quantity  int        `json:"quantity"`
	OrderId   int        `json:"orderId"`
	Price     float64    `json:"price"`
	CreatedAt time.Time  `json:"create_at"`
	UpdatedAt time.Time  `json:"update_at"`
	DeletedAt *time.Time `json:"delete_at,omitempty"`
}

type OrderState int

const (
	Processing OrderState = iota
	ReadyToPay
	Paid
	Completed
	Cancelled
)
