package service

import "context"

type SequencerService struct {
	Repos SequenecerRepository
}

func (svc SequencerService) GetOrderId(ctx context.Context) (orderId int, err error) {
	orderId, err = svc.Repos.GetOrderId(ctx)
	return orderId, err
}
