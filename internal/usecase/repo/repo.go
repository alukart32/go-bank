package repo

import (
	"context"
	"database/sql"
	"fmt"

	"alukart32.com/bank/internal/usecase/repo/db"
)

type SQLRepo struct {
	db *sql.DB
}

func (r *SQLRepo) execTx(ctx context.Context, opts *sql.TxOptions, fn func(q *db.Queries) error) error {
	errCh := make(chan error)
	go func() {
		tx, err := r.db.BeginTx(ctx, opts)
		if err != nil {
			errCh <- err
		}

		// new queries for tx
		qtx := db.New(tx)
		err = fn(qtx)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				errCh <- fmt.Errorf("tx err: %v, rollback err: %v", err, rbErr)
			}
			errCh <- err
		}
		errCh <- tx.Commit()
	}()
	return <-errCh
}
