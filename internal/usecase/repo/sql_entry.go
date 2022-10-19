package repo

import (
	"context"
	"database/sql"

	"alukart32.com/bank/entity"
	"alukart32.com/bank/internal/usecase/repo/db"
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

func (r *EntrySQLRepo) Create(ctx context.Context, e *entity.Entry) (*entity.Entry, error) {
	var result entity.Entry

	err := r.execTx(ctx, nil, func(q *db.Queries) error {
		e, err := q.CreateEntry(ctx, db.CreateEntryParams{
			AccountID: e.AccountID,
			Amount:    e.Amount,
		})
		if err != nil {
			return err
		}

		result = entity.Entry(e)
		return nil
	})

	return &result, err
}

func (r *EntrySQLRepo) Get(ctx context.Context, id int64) (*entity.Entry, error) {
	var result entity.Entry

	err := r.execTx(ctx, nil, func(q *db.Queries) error {
		e, err := q.GetEntry(ctx, id)
		if err != nil {
			return err
		}

		result = entity.Entry(e)
		return nil
	})

	return &result, err
}

func (r *EntrySQLRepo) Update(ctx context.Context, e *entity.Entry) error {
	return r.execTx(ctx, nil, func(q *db.Queries) error {
		return q.UpdateEntry(ctx, db.UpdateEntryParams{
			ID:     e.ID,
			Amount: e.Amount,
		})
	})
}

func (r *EntrySQLRepo) List(ctx context.Context, accountId uuid.UUID) (*[]entity.Entry, error) {
	var result []entity.Entry

	err := r.execTx(ctx, nil, func(q *db.Queries) error {
		entries, err := q.ListEntriesByAccount(ctx, db.ListEntriesByAccountParams{
			AccountID: accountId,
		})
		if err != nil {
			return err
		}

		for _, v := range entries {
			tmp := entity.Entry(v)
			result = append(result, tmp)
		}
		return nil
	})

	return &result, err
}

func (r *EntrySQLRepo) Delete(ctx context.Context, id int64) error {
	return r.execTx(ctx, nil, func(q *db.Queries) error {
		return q.DeleteEntry(ctx, id)
	})
}
