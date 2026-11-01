package model

import "errors"

const (
	DefaultLimit = 10
	MaxLimit     = 50
)

const (
	ProductTypeStandard  = "standard"
	ProductTypeComposite = "composite"
)

var (
	ErrNotFound     = errors.New("record not found")
	ErrInvalidInput = errors.New("invalid input")
)
