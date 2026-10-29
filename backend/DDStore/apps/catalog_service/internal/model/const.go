package model

import "errors"

const MaxLimit = 50
const ProductTableName = "products"

var ErrNotFound = errors.New("record not found")

var ErrProductIsInvalid = errors.New("Product is invalid")

var (
	ErrCategoryIsInvalid = errors.New("category name is required")
	ErrCategoryInUse     = errors.New("category still has products")
	ErrInvalidState      = errors.New("invalid product state")
)

const (
	ProductTypeStandard  = "standard"
	ProductTypeComposite = "composite"
)
