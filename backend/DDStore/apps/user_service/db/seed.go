package db

import (
	"errors"
	"log"

	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/model"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) error {
	p := model.User{}
	if !errors.Is(db.First(&p).Error, gorm.ErrRecordNotFound) {
		log.Println("Data is seeded")
		return nil
	}

	user := []model.User{
		{Name: "system", Email: "system@system.com", Role: "system"},
		{Name: "admin1", Email: "admin@admin.com", Role: "admin"},
	}

	pw := []string{"systempw", "admin01"}

	for i := range user {
		user[i].SetHashPassword([]byte(pw[i]))
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}
