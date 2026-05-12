package repos

import (
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/service"

	"gorm.io/gorm"
)

type mysqlOrderRepo struct {
	DB *gorm.DB
}

func NewMysqlOrderRepo(db *gorm.DB) service.OrderRepository {
	return mysqlOrderRepo{
		DB: db,
	}
}

func (m mysqlOrderRepo) CreateOrder(data *model.Order) (int, error) {
	result := m.DB.Create(data)

	return data.Id, result.Error
}

func (m mysqlOrderRepo) GetOrderById(id int) (data model.Order, err error) {
	result := m.DB.Preload("Items").First(&data, id)

	return data, result.Error
}

func (m mysqlOrderRepo) UpdateOrder(data model.Order) error {
	result := m.DB.Updates(data)
	if result.RowsAffected == 0 && result.Error == nil {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (m mysqlOrderRepo) DeleteOrder(id int) error {
	result := m.DB.Delete(&model.Order{Id: id})
	if result.RowsAffected == 0 && result.Error == nil {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (m mysqlOrderRepo) GetOrders(p model.Paging) (orders []model.Order, err error) {
	result := m.DB.Limit(p.Limit).Offset(p.Offset).Find(&orders)

	return orders, result.Error
}

func (m mysqlOrderRepo) GetHistoryOrders(id int, p model.Paging) (orders []model.Order, err error) {
	result := m.DB.Order("created_at DESC").Limit(p.Limit).Offset(p.Offset).Where("userId", id).Find(&orders)

	return orders, result.Error
}
