package model

import (
	"time"
)

// Category is a document of the categories collection. Ids stay ints (not ObjectIDs)
// because products, the search index and URLs already use int category ids;
// new ids come from the counters collection.
type Category struct {
	Id        int        `bson:"_id"        json:"id"`
	Name      string     `bson:"name"       json:"name"`
	CreatedAt time.Time  `bson:"created_at" json:"create_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"update_at"`
	DeletedAt *time.Time `bson:"deleted_at" json:"delete_at,omitempty"`
}

// ProductCategory is the copy of a category embedded in each product.
type ProductCategory struct {
	Id   int    `bson:"id"   json:"id"`
	Name string `bson:"name" json:"name"`
}
