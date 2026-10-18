package model

import "time"

type Payment struct {
	Id        int       `gorm:"primaryKey"      json:"id"`
	OrderId   int       `gorm:"orderId"      json:"orderId"`
	CreatedAt time.Time `gorm:"autoCreateTime"  json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"  json:"updatedAt"`
}
