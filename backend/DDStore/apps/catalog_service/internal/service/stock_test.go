package service_test

import (
	"log"
	"testing"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestGetStock(t *testing.T) {
	test := []struct {
		productStock          model.Product
		orderItem             service.OrderItemsDto
		quantityStockExpected int
	}{
		{
			productStock:          model.Product{Id: 1, Name: "Milo", Quantity: 10},
			orderItem:             service.OrderItemsDto{ProductId: 1, Quantity: 6},
			quantityStockExpected: 4,
		},
		{
			productStock:          model.Product{Id: 2, Name: "Milk", Quantity: 50},
			orderItem:             service.OrderItemsDto{ProductId: 2, Quantity: 19},
			quantityStockExpected: 31,
		},
	}

	var mockRepos service.MockRepos
	mockRepos.Init()

	log.Println("Create product in mock")
	for _, v := range test {
		mockRepos.CreateProduct(v.productStock)
	}

	order := []service.OrderItemsDto{}
	for _, v := range test {
		order = append(order, v.orderItem)
	}
	log.Println("Call GetStock")

	service.GetStock(order, mockRepos)

	for i := range test {
		log.Println("Get product in mock, id: ", i)

		test[i].productStock, _ = mockRepos.GetProductById(test[i].productStock.Id)
		assert.Equal(t, test[i].quantityStockExpected, test[i].productStock.Quantity)
	}
}

func TestRestoreStock(t *testing.T) {
	test := []struct {
		productStock          model.Product
		orderItem             service.OrderItemsDto
		quantityStockExpected int
	}{
		{
			productStock:          model.Product{Id: 1, Name: "Milo", Quantity: 10},
			orderItem:             service.OrderItemsDto{ProductId: 1, Quantity: 6},
			quantityStockExpected: 16,
		},
		{
			productStock:          model.Product{Id: 2, Name: "Milk", Quantity: 50},
			orderItem:             service.OrderItemsDto{ProductId: 2, Quantity: 19},
			quantityStockExpected: 69,
		},
	}

	var mockRepos service.MockRepos
	mockRepos.Init()
	log.Println("Create product in mock")
	for _, v := range test {
		mockRepos.CreateProduct(v.productStock)
	}

	order := []service.OrderItemsDto{}
	for _, v := range test {
		order = append(order, v.orderItem)
	}
	log.Println("Call RestoreStock")
	service.RestoreStock(order, mockRepos)

	for i := range test {
		log.Println("Get product in mock, id: ", i)

		test[i].productStock, _ = mockRepos.GetProductById(test[i].productStock.Id)
		assert.Equal(t, test[i].quantityStockExpected, test[i].productStock.Quantity)
	}
}
