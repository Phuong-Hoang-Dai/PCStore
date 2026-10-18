package model

import "time"

type Stock struct {
	StockId   int       `json:"stock_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	ProductId int       `json:"product_id"`
}

type StockReservation struct {
	ProductId int       `json:"product_id"`
	OrderId   int       `json:"order_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	StockId   int       `json:"stock_id"`
}
