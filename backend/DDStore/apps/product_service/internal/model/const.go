package model

import "errors"

const MaxLimit = 50
const ProductTableName = "products"

var (
	ErrOutOfStock       = errors.New("some product is out of stock")
	ErrOrderUnvalidated = errors.New("some product is unvalidated")
)
