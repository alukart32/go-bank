package usecase

import (
	"context"
	"errors"

	"alukart32.com/bank/entity"
	"github.com/google/uuid"
)

type entryService struct {
	db EntryRepo
}

func NewEntryService(r EntryRepo) EntryService {
	return &entryService{
		db: r,
	}
}

func (s *entryService) Create(ctx context.Context, e *entity.Entry) (*entity.Entry, error) {
	return nil, errors.New("not implemented yet")
}

func (s *entryService) Get(ctx context.Context, id int64) (*entity.Entry, error) {
	return nil, errors.New("not implemented yet")
}

func (s *entryService) Update(ctx context.Context, e *entity.Entry) (*entity.Entry, error) {
	return nil, errors.New("not implemented yet")
}

func (s *entryService) List(ctx context.Context, accountId uuid.UUID) ([]*entity.Entry, error) {
	return nil, errors.New("not implemented yet")
}

func (s *entryService) Delete(ctx context.Context, id int64) error {
	return errors.New("not implemented yet")
}
