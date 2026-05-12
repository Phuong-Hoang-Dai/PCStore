package service

type ProductService interface {
	GetStock(items []OrderDTO) error
	RestoreStock(items []OrderDTO) error
	GetPriceProduct(items *[]OrderDTO) error
}
