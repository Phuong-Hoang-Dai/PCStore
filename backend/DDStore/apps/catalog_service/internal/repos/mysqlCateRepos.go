package repos

import (
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/service"

	"gorm.io/gorm"
)

type mysqlCateRepo struct {
	DB *gorm.DB
}

func NewMysqlCateRepo(db *gorm.DB) service.CateRepos {
	return mysqlCateRepo{
		DB: db,
	}
}

func (m mysqlCateRepo) CreateCate(data model.Category) (int, error) {
	result := m.DB.Create(&data)

	return data.Id, result.Error
}

func (m mysqlCateRepo) GetCateById(id int) (model.Category, error) {
	data := model.Category{Id: id}
	result := m.DB.First(&data)

	return data, result.Error
}

func (m mysqlCateRepo) UpdateCate(data model.Category) error {
	result := m.DB.Updates(data)

	if result.RowsAffected == 0 && result.Error == nil {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

func (m mysqlCateRepo) DeleteCate(id int) error {
	result := m.DB.Delete(&model.Category{Id: id})
	if result.RowsAffected == 0 && result.Error == nil {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

func (m mysqlCateRepo) GetCates() ([]model.Category, error) {
	data := []model.Category{}
	result := m.DB.Find(&data)

	return data, result.Error
}
