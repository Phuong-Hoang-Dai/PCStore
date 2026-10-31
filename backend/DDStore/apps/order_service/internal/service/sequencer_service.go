package service

import (
	"context"

	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/logger"
)

type SequencerService struct {
	Repos SequenecerRepository
}

func (svc SequencerService) GetOrderId(ctx context.Context) (orderId int, err error) {
	defer logger.Trace(ctx, "SequencerService.GetOrderId")(&orderId, &err)
	orderId, err = svc.Repos.GetOrderId(ctx)
	return orderId, err
}
