package repo

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"alukart32.com/bank/entity"
	"alukart32.com/bank/internal/usecase"
	"alukart32.com/bank/internal/usecase/repo/db"
)

type TransferSQLRepo struct {
	SQLRepo
	mux sync.Mutex
}

func NewTransferSQLRepo(db *sql.DB) *TransferSQLRepo {
	return &TransferSQLRepo{
		SQLRepo: SQLRepo{
			db: db,
		},
	}
}

func (r *TransferSQLRepo) Create(ctx context.Context, transfer *entity.Transfer) (*entity.TransferRes, error) {
	var result entity.TransferRes

	err := r.execTx(ctx, &sql.TxOptions{}, func(q *db.Queries) error {
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

		t, err := q.CreateTransfer(ctx, db.CreateTransferParams{
			FromAccountID: transfer.FromAccountID,
			ToAccountID:   transfer.ToAccountID,
			FromEntryID:   fromEntry.ID,
			ToEntryID:     toEntry.ID,
			Amount:        transfer.Amount,
		})
		if err != nil {
			return err
		}
		result.Transfer = entity.Transfer{
			ID:            t.ID,
			FromAccountID: t.FromAccountID,
			ToAccountID:   t.ToAccountID,
			Amount:        t.Amount,
			FromEntryID:   t.FromEntryID,
			ToEntryID:     t.ToEntryID,
			CreatedAt:     t.CreatedAt,
		}

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
	return &result, err
}

func (r *TransferSQLRepo) Get(ctx context.Context, id int64) (*entity.Transfer, error) {
	var result entity.Transfer

	err := r.execTx(ctx, nil, func(q *db.Queries) error {
		t, err := q.GetTransfer(ctx, id)
		if err != nil {
			return err
		}

		result = entity.Transfer{
			ID:            t.ID,
			FromAccountID: t.FromAccountID,
			ToAccountID:   t.ToAccountID,
			Amount:        t.Amount,
			FromEntryID:   t.FromEntryID,
			ToEntryID:     t.ToEntryID,
			CreatedAt:     t.CreatedAt,
		}
		return nil
	})
	return &result, err
}

func (r *TransferSQLRepo) List(ctx context.Context, params usecase.ListTransferParams) ([]entity.Transfer, error) {
	var result []entity.Transfer

	err := r.execTx(ctx, nil, func(q *db.Queries) error {
		var err error
		switch params.Mode {
		case usecase.ListFromAccount:
			var transfers []db.ListTransfersByFromAccountRow
			transfers, err = q.ListTransfersByFromAccount(ctx, db.ListTransfersByFromAccountParams{
				FromAccountID: params.FromAccountId,
				Limit:         params.Limit,
				Offset:        params.Offset,
			})
			if err != nil {
				return err
			}

			for _, v := range transfers {
				tmp := entity.Transfer(v)
				result = append(result, tmp)
			}
		case usecase.ListToAccount:
			var transfers []db.ListTransfersByToAccountRow
			transfers, err = q.ListTransfersByToAccount(ctx, db.ListTransfersByToAccountParams{
				ToAccountID: params.FromAccountId,
				Limit:       params.Limit,
				Offset:      params.Offset,
			})
			if err != nil {
				return err
			}

			for _, v := range transfers {
				tmp := entity.Transfer(v)
				result = append(result, tmp)
			}
		case usecase.ListByAccounts:
			var transfers []db.ListTransfersByAccountsRow
			transfers, err = q.ListTransfersByAccounts(ctx, db.ListTransfersByAccountsParams{
				ToAccountID:   params.ToAccountId,
				FromAccountID: params.FromAccountId,
				Limit:         params.Limit,
				Offset:        params.Offset,
			})
			if err != nil {
				return err
			}

			for _, v := range transfers {
				tmp := entity.Transfer(v)
				result = append(result, tmp)
			}
		default:
			return errors.New("unsupported list transfer mode")
		}
		return nil
	})

	return result, err
}

func (r *TransferSQLRepo) Rollback(ctx context.Context, id int64) error {
	return r.execTx(ctx, nil, func(q *db.Queries) error {
		// get transfer, fromEntry, toEntry
		transfer, err := q.GetTransfer(ctx, id)
		if err != nil {
			return err
		}
		fromEntry, err := q.GetEntry(ctx, transfer.FromEntryID)
		if err != nil {
			return err
		}
		toEntry, err := q.GetEntry(ctx, transfer.ToEntryID)
		if err != nil {
			return err
		}

		// update fromAccount, toAccount
		r.mux.Lock()
		_, err = q.AddAccountBalance(ctx, db.AddAccountBalanceParams{
			ID:     transfer.FromAccountID,
			Amount: transfer.Amount,
		})
		if err != nil {
			return err
		}

		_, err = q.AddAccountBalance(ctx, db.AddAccountBalanceParams{
			ID:     transfer.ToAccountID,
			Amount: -transfer.Amount,
		})
		if err != nil {
			return err
		}

		// delete fromEntry, toEntry, transfer
		err = q.DeleteTransfer(ctx, transfer.ID)
		if err != nil {
			return err
		}

		err = q.DeleteEntry(ctx, fromEntry.ID)
		if err != nil {
			return err
		}

		err = q.DeleteEntry(ctx, toEntry.ID)
		if err != nil {
			return err
		}

		r.mux.Unlock()
		return nil
	})
}
