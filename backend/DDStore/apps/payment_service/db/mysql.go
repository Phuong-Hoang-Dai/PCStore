package db

import (
	"github.com/Phuong-Hoang-Dai/DDStore/app/payment_service/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/app/payment_service/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDB() (db *gorm.DB, err error) {
	dsn := configs.Cfg.ConnectStr
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&model.Payment{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
