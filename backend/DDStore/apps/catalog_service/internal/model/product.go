package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	Id        int
	Name      string
	Desc      string
	Brand     string
	Cate      Category
	Type      string
	Images    []Media
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type Media struct {
	Url  string
	Role string
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
