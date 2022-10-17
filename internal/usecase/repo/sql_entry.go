package repo

import (
	"context"
	"database/sql"
	"errors"

	"alukart32.com/bank/entity"
	"github.com/google/uuid"
)

type EntrySQLRepo struct {
	SQLRepo
}

func NewEntrySQLRepo(db *sql.DB) *EntrySQLRepo {
	return &EntrySQLRepo{
		SQLRepo: SQLRepo{
			db: db,
		},
	}
}

func (r *EntrySQLRepo) Create(ctx context.Context, e *entity.Entry) error {
	return errors.New("not implemented yet")
}

func (r *EntrySQLRepo) Get(ctx context.Context, id int64) (*entity.Entry, error) {
	return nil, errors.New("not implemented yet")
}

func (r *EntrySQLRepo) Update(ctx context.Context, e *entity.Entry) (*entity.Entry, error) {
	return nil, errors.New("not implemented yet")
}

func (r *EntrySQLRepo) List(ctx context.Context, accountId uuid.UUID) ([]*entity.Entry, error) {
	return nil, errors.New("not implemented yet")
}

func (r *EntrySQLRepo) Delete(id int64) error {
	return errors.New("not implemented yet")
}
