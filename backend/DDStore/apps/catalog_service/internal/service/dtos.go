package service

type OrderItemsDto struct {
	ProductId int     `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
