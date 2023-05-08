package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTransaction(t *testing.T) {
	store := NewStore(testDB)

	account := createRandomAccount(t)
	account1 := createRandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTransactionResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTransaction(context.Background(), TransferTransactionParams{
				FromAccountID: account.ID,
				ToAccountID:   account1.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account.ID, transfer.FromAccountID)
		require.Equal(t, account1.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check Entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account1.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check Accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account1.ID, toAccount.ID)

		// Check Accounts Balance
		diff := account.Balance - fromAccount.Balance
		diff1 := toAccount.Balance - account1.Balance
		require.Equal(t, diff, diff1)
		require.True(t, diff > 0)
		require.True(t, diff%amount == 0)

		k := int(diff / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// Check the final updated balances
	updatedAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	fmt.Println(">> After:", updatedAccount.Balance, updatedAccount1.Balance)
	require.Equal(t, account.Balance-int64(n)*amount, updatedAccount.Balance)
	require.Equal(t, account1.Balance+int64(n)*amount, updatedAccount1.Balance)
}

func TestTransferTransactionDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account := createRandomAccount(t)
	account1 := createRandomAccount(t)

	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account.ID
		toAccountID := account1.ID

		if i%2 == 1 {
			fromAccountID = account1.ID
			toAccountID = account.ID
		}

		go func() {
			_, err := store.TransferTransaction(context.Background(), TransferTransactionParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// Check the final updated balances
	updatedAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	fmt.Println(">> After:", updatedAccount.Balance, updatedAccount1.Balance)
	require.Equal(t, account.Balance, updatedAccount.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
}
