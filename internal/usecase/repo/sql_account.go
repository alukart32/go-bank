package repo

import (
	"context"
	"database/sql"
	"errors"

	"alukart32.com/bank/entity"
	"alukart32.com/bank/internal/usecase/repo/db"
	"github.com/google/uuid"
)

type AccountSQLRepo struct {
	SQLRepo
}

func NewAccountSQLRepo(db *sql.DB) *AccountSQLRepo {
	return &AccountSQLRepo{
		SQLRepo: SQLRepo{
			db: db,
		},
	}
}

func (r *AccountSQLRepo) Create(ctx context.Context, account *entity.Account) (*entity.Account, error) {
	var result entity.Account

	err := r.execTx(ctx, &sql.TxOptions{}, func(q *db.Queries) error {
		a, err := q.CreateAccount(ctx, db.CreateAccountParams{
			ID:       account.ID,
			Owner:    account.Owner,
			Balance:  account.Balance,
			Currency: db.Currency(account.Currency),
		})
		if err != nil {
			return err
		}

		result = entity.Account{
			ID:        a.ID,
			Owner:     a.Owner,
			Balance:   a.Balance,
			Currency:  entity.Currency(a.Currency),
			CreatedAt: a.CreatedAt,
		}
		return nil
	})

	return &result, err
}

func (r *AccountSQLRepo) Get(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	var result entity.Account

	err := r.execTx(ctx, &sql.TxOptions{}, func(q *db.Queries) error {
		a, err := q.GetAccount(ctx, id)
		if err != nil {
			return err
		}

		result = entity.Account{
			ID:        a.ID,
			Owner:     a.Owner,
			Balance:   a.Balance,
			Currency:  entity.Currency(a.Currency),
			CreatedAt: a.CreatedAt,
		}
		return nil
	})

	return &result, err
}

func (r *AccountSQLRepo) Update(ctx context.Context, a *entity.Account) (*entity.Account, error) {
	return nil, errors.New("not implemented yet")
}

func (r *AccountSQLRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.New("not implemented yet")
}
