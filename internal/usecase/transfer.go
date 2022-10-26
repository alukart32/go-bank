package usecase

import (
	"context"
	"errors"

	"alukart32.com/bank/entity"
	"alukart32.com/bank/pkg/zerologx"
)

type transferService struct {
	db TransferRepo
	l  zerologx.Logger
}

func NewTransferService(r TransferRepo, l zerologx.Logger) TransferService {
	return &transferService{
		db: r,
		l:  l,
	}
}

func (s *transferService) Transfer(ctx context.Context, t entity.Transfer) (entity.TransferRes, error) {
	return entity.TransferRes{}, errors.New("not implemented yet")
}

func (s *transferService) Get(ctx context.Context, id int64) (entity.Transfer, error) {
	return entity.Transfer{}, errors.New("not implemented yet")
}

func (s *transferService) List(ctx context.Context, params ListTransferParams) ([]entity.Transfer, error) {
	return nil, errors.New("not implemented yet")
}

func (s *transferService) Rollback(ctx context.Context, id int64) error {
	return errors.New("not implemented yet")
}
