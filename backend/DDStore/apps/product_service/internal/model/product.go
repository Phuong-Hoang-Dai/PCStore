package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	Id         int            `json:"id" gorm:"primaryKey;autoIncrement;<-:create"`
	Name       string         `json:"name" gorm:"column:name"`
	Desc       string         `json:"description" gorm:"column:description"`
	Price      float64        `json:"price" gorm:"column:price"`
	Quantity   int            `json:"quantity" gorm:"column:quantity"`
	Image      string         `json:"image" gorm:"column:string"`
	CategoryID int            `json:"category_id" gorm:"column:category_id"`
	Cate       Category       `json:"category" gorm:"foreignKey:CategoryID;references:Id"`
	CreatedAt  time.Time      `json:"create_at"`
	UpdatedAt  time.Time      `json:"update_at"`
	DeletedAt  gorm.DeletedAt `json:"delete_at" gorm:"index"`
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
