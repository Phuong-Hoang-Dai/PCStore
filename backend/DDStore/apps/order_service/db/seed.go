package db

import (
	"errors"
	"log"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) error {
	o := model.Order{}

	if !errors.Is(db.First(&o).Error, gorm.ErrRecordNotFound) {
		log.Println("Data is seeded")
		return nil
	}

	orders := []model.Order{
		{State: 0, Total: 0},
		{State: 0, Total: 0},
	}

	if err := db.Create(&orders).Error; err != nil {
		return err
	}

	orderItems := []model.OrderItem{
		{ProductId: 1, Quantity: 2, OrderId: orders[0].Id},
		{ProductId: 2, Quantity: 1, OrderId: orders[0].Id},

		{ProductId: 3, Quantity: 1, OrderId: orders[1].Id},
		{ProductId: 4, Quantity: 2, OrderId: orders[1].Id},
	}

	if err := db.Create(&orderItems).Error; err != nil {
		return err
	}

	return nil
}
