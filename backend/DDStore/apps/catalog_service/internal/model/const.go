package model

import "errors"

const MaxLimit = 50
const ProductTableName = "products"

var ErrNotFound = errors.New("record not found")

var ErrProductIsInvalid = errors.New("Product is invalid")

const (
	ProductTypeStandard  = "standard"
	ProductTypeComposite = "composite"
)
