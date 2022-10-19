package repo

import (
	"context"
	"database/sql"

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

	err := r.execTx(ctx, nil, func(q *db.Queries) error {
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

func (r *AccountSQLRepo) UpdateOwner(ctx context.Context, id uuid.UUID, owner string) (*entity.Account, error) {
	var result entity.Account

	err := r.execTx(ctx, nil, func(q *db.Queries) error {
		a, err := q.UpdateAccountOwner(ctx, db.UpdateAccountOwnerParams{
			ID:    id,
			Owner: owner,
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

func (r *AccountSQLRepo) UpdateBalance(ctx context.Context, id uuid.UUID, amount int64) (*entity.Account, error) {
	var result entity.Account

	err := r.execTx(ctx, nil, func(q *db.Queries) error {
		a, err := q.AddAccountBalance(ctx, db.AddAccountBalanceParams{
			ID:     id,
			Amount: amount,
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

func (r *AccountSQLRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.execTx(ctx, nil, func(q *db.Queries) error {
		return q.DeleteAccount(ctx, id)
	})
}
