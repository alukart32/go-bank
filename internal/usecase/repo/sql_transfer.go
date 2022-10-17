package repo

import (
	"context"
	"database/sql"
	"errors"

	"alukart32.com/bank/entity"
	"alukart32.com/bank/internal/usecase/repo/db"
	"github.com/google/uuid"
)

type TransferSQLRepo struct {
	SQLRepo
}

func NewTransferSQLRepo(db *sql.DB) *TransferSQLRepo {
	return &TransferSQLRepo{
		SQLRepo: SQLRepo{
			db: db,
		},
	}
}

func (r *TransferSQLRepo) Create(ctx context.Context, transfer *entity.Transfer) (*entity.TransferRes, error) {
	result := &entity.TransferRes{}

	err := r.execTx(ctx, &sql.TxOptions{}, func(q *db.Queries) error {
		t, err := q.CreateTransfer(ctx, db.CreateTransferParams{
			FromAccountID: transfer.FromAccountID,
			ToAccountID:   transfer.ToAccountID,
			Amount:        transfer.Amount,
		})
		if err != nil {
			return err
		}
		result.Transfer = entity.Transfer(t)

		fromEntry, err := q.CreateEntry(ctx, db.CreateEntryParams{
			AccountID: transfer.FromAccountID,
			Amount:    -transfer.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntry = entity.Entry(fromEntry)

		toEntry, err := q.CreateEntry(ctx, db.CreateEntryParams{
			AccountID: transfer.ToAccountID,
			Amount:    transfer.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry = entity.Entry(toEntry)

		// update accounts
		r.mux.Lock()
		fromAccount, err := q.AddAccountBalance(ctx, db.AddAccountBalanceParams{
			ID:     transfer.FromAccountID,
			Amount: -transfer.Amount,
		})
		if err != nil {
			return err
		}
		result.FromAccount = entity.Account{
			ID:        fromAccount.ID,
			Owner:     fromAccount.Owner,
			Balance:   fromAccount.Balance,
			Currency:  entity.Currency(fromAccount.Currency),
			CreatedAt: fromAccount.CreatedAt,
		}

		toAccount, err := q.AddAccountBalance(ctx, db.AddAccountBalanceParams{
			ID:     transfer.ToAccountID,
			Amount: transfer.Amount,
		})
		if err != nil {
			return err
		}
		result.ToAccount = entity.Account{
			ID:        toAccount.ID,
			Owner:     toAccount.Owner,
			Balance:   toAccount.Balance,
			Currency:  entity.Currency(toAccount.Currency),
			CreatedAt: toAccount.CreatedAt,
		}

		r.mux.Unlock()
		return nil
	})
	return result, err
}

func (r *TransferSQLRepo) Get(ctx context.Context, id int64) (*entity.Transfer, error) {
	return nil, errors.New("not implemented yet")
}
func (r *TransferSQLRepo) List(ctx context.Context, accountId uuid.UUID) ([]*entity.Transfer, error) {
	return nil, errors.New("not implemented yet")
}
func (r *TransferSQLRepo) Delete(ctx context.Context, id int64, accountId uuid.UUID) error {
	return errors.New("not implemented yet")
}
