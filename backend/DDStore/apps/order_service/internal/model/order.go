package model

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
	DeletedAt time.Time   `json:"delete_at" `
}

type OrderItem struct {
	ProductId string    `json:"productId"`
	Quantity  int       `json:"quantity"`
	OrderId   int       `json:"orderId"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
	DeletedAt time.Time `json:"delete_at" `

	Order Order
}

type OrderState int

const (
	Processing OrderState = iota
	ReadyToPay
	Payed
	Completed
	Cancelled
)

func (order Order) Validate() bool {
	if order.Items == nil {
		return false
	}

	return order.Total != 0
}

func (o *Order) CalcTotal() {
	for _, v := range o.Items {
		o.Total += v.Price * float64(v.Quantity)
	}
}

type Paging struct {
	Limit  int
	Offset int
}

func (p *Paging) Process() {
	if p.Limit > MaxLimit {
		p.Limit = MaxLimit
	}
	if p.Limit < 0 {
		p.Limit = 0
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}
