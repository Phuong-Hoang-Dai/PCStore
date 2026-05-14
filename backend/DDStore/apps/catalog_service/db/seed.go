package db

import (
	"errors"
	"log"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) error {

	p := model.Category{}
	if !errors.Is(db.First(&p).Error, gorm.ErrRecordNotFound) {
		log.Println("Data is seeded")
		return nil
	}

	cates := []model.Category{
		{Id: 1, Name: "PCGAMING"},
		{Id: 2, Name: "PCDESIGN"},
		{Id: 3, Name: "PCWORK"},
	}

	if err := db.Create(&cates).Error; err != nil {
		return err
	}

	products := []model.Product{
		{Name: "PC BEST FOR GAMING", Desc: "i5 11400F- GTX 1660 Super 6GB"},
		{Name: "PC CHƠI GAME HIỆU SUẤT CAO", Desc: "RTX 3060 12GB - 12400F "},
		{Name: "PC DESIGNER -3D RENDER - EDIT VIDEO", Desc: "i5 14700KF - RTX 5050 8GB"},
		{Name: "PC DESIGNER EDIT VIDEO", Desc: "MD Ryzen 9 9950X3D - RTX 5090 32GB OC"},
		{Name: "PC HOME OFFICE", Desc: "Core i3 10105 - RAM 8GB- SSD 256GB- Kèm Màn Hình"},
		{Name: "PC HOME OFFICE", Desc: "Core i3 10105 - RAM 8GB- SSD 256GB- Kèm Màn Hình"},
		{Name: "PC HOME OFFICE", Desc: "Core i3 10105 - RAM 8GB- SSD 256GB- Kèm Màn Hình"},
		{Name: "PC HOME OFFICE", Desc: "Core i3 10105 - RAM 8GB- SSD 256GB- Kèm Màn Hình"},
		{Name: "PC HOME OFFICE", Desc: "Core i3 10105 - RAM 8GB- SSD 256GB- Kèm Màn Hình"},
	}

	if err := db.Create(&products).Error; err != nil {
		return err
	}

	return nil
}
