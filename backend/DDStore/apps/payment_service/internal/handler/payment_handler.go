package handler

import (
	"github.com/Phuong-Hoang-Dai/DDStore/app/payment_service/internal/repos"
	"github.com/Phuong-Hoang-Dai/DDStore/app/payment_service/internal/service"
	"gorm.io/gorm"
)

type PaymentService struct {
	paymentRepos service.PaymentRepository
}

func Init(db *gorm.DB) PaymentService {
	p := PaymentService{
		paymentRepos: repos.NewMysqlPaymentRepo(db),
	}

	return p
}
