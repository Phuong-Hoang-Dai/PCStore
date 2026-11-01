package service

type OrderItemsDto struct {
	ProductId string  `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
