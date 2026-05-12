package db

import (
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDB() (db *gorm.DB, err error) {
	dsn := configs.Cfg.ConnectStr
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&model.Order{}, &model.OrderItem{})
	if err != nil {
		return nil, err
	}
	err = SeedData(db)

	if err != nil {
		return db, err
	}

	return db, nil
}
