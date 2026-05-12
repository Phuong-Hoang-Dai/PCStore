package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	Id        int            `json:"orderId" gorm:"column:id"`
	State     int            `json:"state" gorm:"column:state"`
	Items     []OrderItem    `json:"items" gorm:"foreignKey:OrderId;references:Id"`
	Total     float64        `json:"total" gorm:"column:total"`
	UserId    int            `json:"userId" gorm:"column:userId"`
	CreatedAt time.Time      `json:"create_at"`
	UpdatedAt time.Time      `json:"update_at"`
	DeletedAt gorm.DeletedAt `json:"delete_at" gorm:"index"`
}

type OrderItem struct {
	ProductId int            `json:"productId" gorm:"column:productId"`
	Quantity  int            `json:"quantity" gorm:"column:quantity"`
	OrderId   int            `json:"orderId" gorm:"column:orderId"`
	Price     float64        `json:"price" gorm:"column:price"`
	CreatedAt time.Time      `json:"create_at"`
	UpdatedAt time.Time      `json:"update_at"`
	DeletedAt gorm.DeletedAt `json:"delete_at" gorm:"index"`

	Order Order `gorm:"foreignKey:OrderId;references:id"`
}

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
