package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"alukart32.com/bank/config"
	"alukart32.com/bank/pkg/postgres"
	"alukart32.com/bank/pkg/random"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testDB   *sql.DB
	testConf config.Config
)

func TestMain(m *testing.M) {
	var err error
	testConf, err = config.New(config.Test)
	if err != nil {
		log.Fatal("cannot get config: ", err)
	}

	testDB, err = postgres.New(testConf.DB)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	os.Exit(m.Run())
}
func TestCreateAccount(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	createRandomAccount(t, New(tx))

	if err = tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}

func TestGetAccount(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	qtx := New(tx)

	creationTime := time.Now()
	account1 := createRandomAccount(t, qtx)
	account2, err := qtx.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Owner, account2.Owner)
	assert.Equal(t, account1.Balance, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, creationTime, account2.CreatedAt, time.Second)

	if err = tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateAccount(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	qtx := New(tx)

	account := createRandomAccount(t, qtx)

	amount := -random.Int64(1, 2000)
	arg := AddAccountBalanceParams{
		ID:     account.ID,
		Amount: amount,
	}
	updatedAccount, err := qtx.AddAccountBalance(context.Background(), arg)

	require.NoError(t, err)
	assert.Equal(t, account.Balance+amount, updatedAccount.Balance)

	if err = tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteAccount(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	qtx := New(tx)

	account1 := createRandomAccount(t, qtx)

	err = qtx.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := qtx.GetAccount(context.Background(), account1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)

	if err = tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}

func TestListAccounts(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	qtx := New(tx)

	for i := 0; i < 10; i++ {
		createRandomAccount(t, qtx)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 0,
	}

	list, err := qtx.ListAccounts(context.Background(), arg)
	require.NoError(t, err)

	assert.Equal(t, 5, len(list))
	for _, v := range list {
		assert.NotEmpty(t, v)
	}

	if err = tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateEntry(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	qtx := New(tx)

	account := createRandomAccount(t, qtx)

	amount := -random.Int64(1, 2000)
	addAccountBalance := AddAccountBalanceParams{
		ID:     account.ID,
		Amount: amount,
	}
	_, err = qtx.AddAccountBalance(context.Background(), addAccountBalance)
	require.NoError(t, err)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    amount,
	}
	r, err := qtx.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	assert.Equal(t, account.ID, r.AccountID)
	assert.Equal(t, amount, r.Amount)

	if err = tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}

func TestGetEntry(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	qtx := New(tx)

	account := createRandomAccount(t, qtx)

	amount := random.Int64(1, 200000)
	r, err := qtx.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: account.ID,
		Amount:    amount,
	})
	require.NoError(t, err)
	assert.Equal(t, account.ID, r.AccountID)
	assert.Equal(t, amount, r.Amount)

	r, err = qtx.GetEntry(context.Background(), r.ID)
	require.NoError(t, err)
	assert.Equal(t, account.ID, r.AccountID)

	if err = tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}

func TestListEntriesByAccount(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	qtx := New(tx)

	account := createRandomAccount(t, qtx)
	for i := 0; i < 4; i++ {
		time.Sleep(25 * time.Millisecond)

		amount := int64(10)
		_, err := qtx.CreateEntry(context.Background(), CreateEntryParams{
			AccountID: account.ID,
			Amount:    amount,
		})
		require.NoError(t, err)
	}
	list, err := qtx.ListEntriesByAccount(context.Background(), ListEntriesByAccountParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    0,
	})
	require.NoError(t, err)

	for i := 1; i < len(list); i++ {
		if list[i].AccountID != account.ID {
			t.Errorf("get wrong entry; expect accountID %v, actual %v", account.ID, list[i].AccountID)
		}
	}

	if err = tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateEntry(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	qtx := New(tx)

	account := createRandomAccount(t, qtx)

	entry, err := qtx.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: account.ID,
		Amount:    random.Int64(1, 200000),
	})
	require.NoError(t, err)

	newAmmount := random.Int64(1, 1000)
	err = qtx.UpdateEntry(context.Background(), UpdateEntryParams{
		ID:     entry.ID,
		Amount: newAmmount,
	})
	require.NoError(t, err)

	updated, err := qtx.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	assert.Equal(t, newAmmount, updated.Amount)

	if err = tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteEntry(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	qtx := New(tx)

	account := createRandomAccount(t, qtx)

	entry, err := qtx.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: account.ID,
		Amount:    random.Int64(1, 200000),
	})
	require.NoError(t, err)

	err = qtx.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	_, err = qtx.GetEntry(context.Background(), entry.ID)
	require.ErrorIs(t, sql.ErrNoRows, err)

	if err = tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateTransfer(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	qtx := New(tx)

	amount := random.Int64(1, 200000)
	fromAccount := createRandomAccount(t, qtx)
	toAccount := createRandomAccount(t, qtx)
	fromEntry, err := qtx.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: fromAccount.ID,
		Amount:    -amount,
	})
	require.NoError(t, err)
	toEntry, err := qtx.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: toAccount.ID,
		Amount:    amount,
	})
	require.NoError(t, err)

	r, err := qtx.CreateTransfer(context.Background(), CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		FromEntryID:   fromEntry.ID,
		ToEntryID:     toEntry.ID,
		Amount:        amount,
	})
	require.NoError(t, err)
	assert.Equal(t, fromAccount.ID, r.FromAccountID)
	assert.Equal(t, toAccount.ID, r.ToAccountID)
	assert.Equal(t, amount, r.Amount)

	if err = tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}

func TestGetTransfer(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	qtx := New(tx)

	amount := random.Int64(1, 200000)
	fromAccount := createRandomAccount(t, qtx)
	toAccount := createRandomAccount(t, qtx)
	fromEntry, err := qtx.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: fromAccount.ID,
		Amount:    -amount,
	})
	require.NoError(t, err)
	toEntry, err := qtx.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: toAccount.ID,
		Amount:    amount,
	})
	require.NoError(t, err)

	r, err := qtx.CreateTransfer(context.Background(), CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		FromEntryID:   fromEntry.ID,
		ToEntryID:     toEntry.ID,
		Amount:        amount,
	})
	require.NoError(t, err)

	transfer, err := qtx.GetTransfer(context.Background(), r.ID)
	require.NoError(t, err)
	assert.Equal(t, r.ID, transfer.ID)
	assert.Equal(t, r.FromAccountID, transfer.FromAccountID)
	assert.Equal(t, r.ToAccountID, transfer.ToAccountID)

	if err = tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}

func TestListTransfersByFromAccount(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	qtx := New(tx)

	amount := random.Int64(1, 100)
	fromAccount := createRandomAccount(t, qtx)
	toAccounts := make([]Account, 5)
	for i := 0; i < 5; i++ {
		toAccounts[i] = createRandomAccount(t, qtx)

	}

	transfers := make([]Transfer, 5)
	for i, v := range toAccounts {
		fromEntry, err := qtx.CreateEntry(context.Background(), CreateEntryParams{
			AccountID: fromAccount.ID,
			Amount:    -amount,
		})
		require.NoError(t, err)
		toEntry, err := qtx.CreateEntry(context.Background(), CreateEntryParams{
			AccountID: v.ID,
			Amount:    amount,
		})
		require.NoError(t, err)

		transfer, err := qtx.CreateTransfer(context.Background(), CreateTransferParams{
			FromAccountID: fromAccount.ID,
			ToAccountID:   v.ID,
			FromEntryID:   fromEntry.ID,
			ToEntryID:     toEntry.ID,
			Amount:        amount,
		})
		require.NoError(t, err)

		transfers[i] = transfer
	}

	for _, v := range transfers {
		if v.FromAccountID != fromAccount.ID {
			t.Errorf("get wrong transfer; expect fromAccount.ID: %v, actual %v", fromAccount.ID, v.FromAccountID)
		}
	}

	if err = tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}

func TestListTransfersByToAccount(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	qtx := New(tx)

	amount := random.Int64(1, 100)
	toAccount := createRandomAccount(t, qtx)
	fromAccounts := make([]Account, 5)
	for i := 0; i < 5; i++ {
		fromAccounts[i] = createRandomAccount(t, qtx)

	}

	transfers := make([]Transfer, 5)
	for i, v := range fromAccounts {
		fromEntry, err := qtx.CreateEntry(context.Background(), CreateEntryParams{
			AccountID: v.ID,
			Amount:    -amount,
		})
		require.NoError(t, err)
		toEntry, err := qtx.CreateEntry(context.Background(), CreateEntryParams{
			AccountID: toAccount.ID,
			Amount:    amount,
		})
		require.NoError(t, err)

		transfer, err := qtx.CreateTransfer(context.Background(), CreateTransferParams{
			FromAccountID: v.ID,
			ToAccountID:   toAccount.ID,
			FromEntryID:   fromEntry.ID,
			ToEntryID:     toEntry.ID,
			Amount:        amount,
		})
		require.NoError(t, err)

		transfers[i] = transfer
	}

	for _, v := range transfers {
		if v.ToAccountID != toAccount.ID {
			t.Errorf("get wrong transfer; expect toAccount.ID: %v, actual %v", toAccount.ID, v.FromAccountID)
		}
	}

	if err = tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteTransfer(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	qtx := New(tx)

	amount := random.Int64(1, 200000)
	fromAccount := createRandomAccount(t, qtx)
	toAccount := createRandomAccount(t, qtx)
	fromEntry, err := qtx.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: fromAccount.ID,
		Amount:    -amount,
	})
	require.NoError(t, err)
	toEntry, err := qtx.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: toAccount.ID,
		Amount:    amount,
	})
	require.NoError(t, err)

	r, err := qtx.CreateTransfer(context.Background(), CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		FromEntryID:   fromEntry.ID,
		ToEntryID:     toEntry.ID,
		Amount:        amount,
	})
	require.NoError(t, err)

	err = qtx.DeleteTransfer(context.Background(), r.ID)
	require.NoError(t, err)

	if err = tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}

// TODO: replace with golden files
func createRandomAccount(t *testing.T, queries *Queries) Account {
	arg := CreateAccountParams{
		ID:       uuid.New(),
		Owner:    string(random.String(20)),
		Balance:  random.Int64(1, 900000),
		Currency: Currency(random.GetString([]string{string(CurrencyRUB), string(CurrencyUSD)}...)),
	}

	acc, err := queries.CreateAccount(context.Background(), arg)
	if err != nil {
		t.Error(err)
	}

	require.NoError(t, err)

	return acc
}
