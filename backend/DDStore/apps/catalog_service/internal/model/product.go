package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Product struct {
	Id        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string        `bson:"name"          json:"name"`
	Desc      string        `bson:"description"   json:"description"`
	Brand     string        `bson:"brand"         json:"brand"`
	Cate      Category      `bson:"cate"          json:"cate"`
	Type      string        `bson:"type"          json:"type"`
	Images    []Media       `bson:"images"        json:"images"`
	CreatedAt time.Time     `bson:"created_at"    json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updated_at"`
	DeletedAt *time.Time    `bson:"deleted_at"    json:"deleted_at,omitempty"`
}

type Media struct {
	Url  string `bson:"url"  json:"url"`
	Role string `bson:"role" json:"role"`
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
