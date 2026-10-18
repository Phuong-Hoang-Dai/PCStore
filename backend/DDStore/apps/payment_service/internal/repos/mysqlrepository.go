package repos

import (
	"github.com/Phuong-Hoang-Dai/DDStore/app/payment_service/internal/service"

	"gorm.io/gorm"
)

type mysqlPaymentRepo struct {
	DB *gorm.DB
}

func NewMysqlPaymentRepo(db *gorm.DB) service.PaymentRepository {
	return mysqlPaymentRepo{
		DB: db,
	}
}
