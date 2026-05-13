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
	Image     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type Option struct {
	Name       string
	Value      []string
	IsRequired bool
}

type OptionValue struct {
	Name  string
	Value string
}

type Variant struct {
	Id     int
	Sku    string
	option OptionValue
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
