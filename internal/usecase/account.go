package usecase

import (
	"context"
	"errors"

	"alukart32.com/bank/entity"
	"alukart32.com/bank/pkg/zerologx"
	"github.com/google/uuid"
)

type accountService struct {
	db AccountRepo
	l  zerologx.Logger
}

func NewAccountService(r AccountRepo, l zerologx.Logger) AccountService {
	return &accountService{
		db: r,
		l:  l,
	}
}

func (s *accountService) Create(ctx context.Context, a entity.Account) (uuid.UUID, error) {
	return [16]byte{}, errors.New("not implemented yet")
}

func (s *accountService) Get(ctx context.Context, id uuid.UUID) (entity.Account, error) {
	return entity.Account{}, errors.New("not implemented yet")
}

func (s *accountService) UpdateOwner(ctx context.Context, id uuid.UUID, owner string) (entity.Account, error) {
	return entity.Account{}, errors.New("not implemented yet")
}

func (s *accountService) AddBalance(ctx context.Context, id uuid.UUID, amount int64) (entity.Account, error) {
	return entity.Account{}, errors.New("not implemented yet")
}

func (s *accountService) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.New("not implemented yet")
}

func (s *accountService) ListEntries(ctx context.Context, id uuid.UUID) ([]entity.Entry, error) {
	return nil, errors.New("not implemented yet")
}

func (s *accountService) ListTransfers(ctx context.Context, id uuid.UUID) ([]entity.Transfer, error) {
	return nil, errors.New("not implemented yet")
}
