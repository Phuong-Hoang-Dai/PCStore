package repos

import (
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/service"

	"gorm.io/gorm"
)

type mysqlProductRepo struct {
	DB *gorm.DB
}

func NewMysqlProductRepo(db *gorm.DB) service.ProductRepos {
	return mysqlProductRepo{
		DB: db,
	}
}

func (m mysqlProductRepo) CreateProduct(data model.Product) (int, error) {
	result := m.DB.Create(&data)

	return data.Id, result.Error
}

func (m mysqlProductRepo) GetProductById(id int) (model.Product, error) {
	data := model.Product{Id: id}
	result := m.DB.First(&data).Preload("Cate")

	return data, result.Error
}

func (m mysqlProductRepo) UpdateProduct(data model.Product) error {
	result := m.DB.Updates(data)

	if result.RowsAffected == 0 && result.Error == nil {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

func (m mysqlProductRepo) UpdateProducts(data []model.Product) error {
	return m.DB.Transaction(func(tx *gorm.DB) error {
		for i := range data {
			result := tx.Updates(data[i])
			if result.RowsAffected == 0 && result.Error == nil {
				return gorm.ErrRecordNotFound
			}
		}
		return nil
	})
}

func (m mysqlProductRepo) DeleteProduct(id int) error {
	result := m.DB.Delete(&model.Product{Id: id})
	if result.RowsAffected == 0 && result.Error == nil {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

func (m mysqlProductRepo) GetProducts(p model.Paging) ([]model.Product, error) {
	data := []model.Product{}
	result := m.DB.Preload("Cate").Limit(p.Limit).Offset(p.Offset).Find(&data)

	return data, result.Error
}

func (m mysqlProductRepo) GetProductsByCate(p model.Paging, cate model.Category) ([]model.Product, error) {
	data := []model.Product{}
	result := m.DB.Preload("Cate").Where("category_id", cate.Id).Limit(p.Limit).Offset(p.Offset).Find(&data)

	return data, result.Error
}
