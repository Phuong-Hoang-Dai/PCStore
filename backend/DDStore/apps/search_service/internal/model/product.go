package model

import "time"

// Product is the document stored in and returned from Elasticsearch.
// It mirrors the catalog_service product model, but uses a string Id
// (the Elasticsearch document _id) and JSON tags only, since Elasticsearch
// speaks JSON.
type Product struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Desc      string    `json:"description"`
	Brand     string    `json:"brand"`
	Cate      Category  `json:"cate"`
	Price     float32   `json:"price"`
	Type      string    `json:"type"`
	Images    []Media   `json:"images"`
	Upgrade   any       `json:"upgrade,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Media struct {
	Url  string `json:"url"`
	Role string `json:"role"`
}

type Upgrade struct {
	Name       string      `json:"name"`
	Components []Component `json:"commponents"`
}

type Component struct {
	Id    string  `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}
