package db

import (
	"avancedGo/util"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwnerName(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccoutn(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {

	randomAccount1 := createRandomAccount(t)
	randomAccount2, err := testQueries.GetAccount(context.Background(), randomAccount1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, randomAccount2)

	require.Equal(t, randomAccount1.ID, randomAccount2.ID)
	require.Equal(t, randomAccount1.Balance, randomAccount2.Balance)
	require.Equal(t, randomAccount1.Currency, randomAccount2.Currency)
	require.Equal(t, randomAccount1.Owner, randomAccount2.Owner)
	require.Equal(t, randomAccount1.Currency, randomAccount2.Currency)

	require.WithinDuration(t, randomAccount1.CreatedAt, randomAccount2.CreatedAt, time.Second)

}
