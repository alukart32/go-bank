package usecase

import (
	"context"
	"errors"

	"alukart32.com/bank/entity"
	"github.com/google/uuid"
)

type accountService struct {
	db AccountRepo
}

func NewAccountService(r AccountRepo) AccountService {
	return &accountService{
		db: r,
	}
}

func (s *accountService) Create(ctx context.Context, a *entity.Account) error {
	return errors.New("not implemented yet")
}

func (s *accountService) Get(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	return nil, errors.New("not implemented yet")
}

func (s *accountService) Update(ctx context.Context, a *entity.Account) (*entity.Account, error) {
	return nil, errors.New("not implemented yet")
}

func (s *accountService) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.New("not implemented yet")
}

func (s *accountService) ListEntries(ctx context.Context, id uuid.UUID) ([]*entity.Entry, error) {
	return nil, errors.New("not implemented yet")
}

func (s *accountService) ListTransfers(ctx context.Context, id uuid.UUID) ([]*entity.Transfer, error) {
	return nil, errors.New("not implemented yet")
}
