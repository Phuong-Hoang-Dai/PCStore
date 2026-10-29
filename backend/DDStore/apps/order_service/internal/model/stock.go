package model

import "time"

type Stock struct {
	ProductId string    `json:"productId"`
	Quantity  int       `json:"quantity"`
	UpdatedAt time.Time `json:"updated_at"`
}
