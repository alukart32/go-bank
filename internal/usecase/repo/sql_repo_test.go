package repo

import (
	"context"
	"fmt"
	"testing"

	"alukart32.com/bank/pkg/random"
	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	sqlRepo := NewSQLRepo(testDB)

	fromAccount := createRandomAccount(t, sqlRepo.Queries)
	toAccount := createRandomAccount(t, sqlRepo.Queries)

	n := 5
	amount := random.Int64(1, 200)
	errors := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			t, err := sqlRepo.Transfer(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})

			errors <- err
			results <- t
		}()
	}

	// check results
	existed := make(map[int]bool)

	// check results
	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// check accounts
		fromAccountTx := result.FromAccount
		require.NotEmpty(t, fromAccountTx)
		require.Equal(t, fromAccount.ID, fromAccountTx.ID)

		toAccountTx := result.ToAccount
		require.NotEmpty(t, toAccountTx)
		require.Equal(t, toAccount.ID, toAccountTx.ID)

		// check balances
		diff1 := fromAccount.Balance - fromAccountTx.Balance
		diff2 := toAccountTx.Balance - toAccount.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)

		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedFromAccount, err := sqlRepo.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedToAccount, err := sqlRepo.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	require.Equal(t, fromAccount.Balance-int64(n)*amount, updatedFromAccount.Balance)
	require.Equal(t, toAccount.Balance+int64(n)*amount, updatedToAccount.Balance)
}

func TestTransferDeadlock(t *testing.T) {
	sqlRepo := NewSQLRepo(testDB)

	fromAccount := createRandomAccount(t, sqlRepo.Queries)
	toAccount := createRandomAccount(t, sqlRepo.Queries)
	fmt.Printf(">> before: [ fromAccount.Balance: %v, toAccount.Balance: %v ]\n", fromAccount.Balance, toAccount.Balance)

	n := 10
	amount := random.Int64(1, 200)
	errors := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID, toAccountID := fromAccount.ID, toAccount.ID

		if i%2 == 1 {
			fromAccountID, toAccountID = toAccountID, fromAccountID
		}

		go func() {
			_, err := sqlRepo.Transfer(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errors <- err
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedFromAccount, err := sqlRepo.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedToAccount, err := sqlRepo.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	require.Equal(t, fromAccount.Balance, updatedFromAccount.Balance)
	require.Equal(t, toAccount.Balance, updatedToAccount.Balance)
}
