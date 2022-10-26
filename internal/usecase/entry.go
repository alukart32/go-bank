package usecase

import (
	"context"
	"errors"

	"alukart32.com/bank/entity"
	"alukart32.com/bank/pkg/zerologx"
	"github.com/google/uuid"
)

type entryService struct {
	db EntryRepo
	l  zerologx.Logger
}

func NewEntryService(r EntryRepo, l zerologx.Logger) EntryService {
	return &entryService{
		db: r,
		l:  l,
	}
}

func (s *entryService) Create(ctx context.Context, e entity.Entry) (entity.Entry, error) {
	return entity.Entry{}, errors.New("not implemented yet")
}

func (s *entryService) Get(ctx context.Context, id int64) (entity.Entry, error) {
	return entity.Entry{}, errors.New("not implemented yet")
}

func (s *entryService) UpdateAmount(ctx context.Context, e entity.Entry) (entity.Entry, error) {
	return entity.Entry{}, errors.New("not implemented yet")
}

func (s *entryService) List(ctx context.Context, accountId uuid.UUID) ([]entity.Entry, error) {
	return nil, errors.New("not implemented yet")
}

func (s *entryService) Delete(ctx context.Context, id int64) error {
	return errors.New("not implemented yet")
}
