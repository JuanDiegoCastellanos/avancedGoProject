package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTX(ctx context.Context, arg TransferTXParams) (TransferTXResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	DB *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		DB:      db,
		Queries: New(db),
	}
}

// Function to execute a transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queries := New(tx)

	if err = fn(queries); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return fmt.Errorf("have been ocurred some errors during transaction and rollback: %v, %v", err, errRollback)
		}
		return err
	}
	return tx.Commit()
}

type TransferTXParams struct {
	FromAccountID int64 `json: "from_account_id"`
	ToAccountID   int64 `json: "to_account_id"`
	Amount        int64 `json: "amount"`
}
type TransferTXResult struct {
	Transfer    Transfer `json: "transfer"`
	FromAccount Account  `json: "from_account"`
	ToAccount   Account  `json: "to_account"`
	FromEntry   Entry    `json: "from_entry"`
	ToEntry     Entry    `json: "to_entry"`
}

//var txKey = struct{}{}

// TransferTX Function to execute a Transfer using a transaction
// TransferTX per forms a money transfer from one account to the other.
// It creates the transfer, add account entries, and update accounts' balance within a single database transaction
func (store *SQLStore) TransferTX(ctx context.Context, arg TransferTXParams) (TransferTXResult, error) {
	var result TransferTXResult

	err := store.execTx(ctx, func(q *Queries) error {
		var errTrans error

		// txName := ctx.Value(txKey)

		// fmt.Println(txName, "create transfer")

		result.Transfer, errTrans = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if errTrans != nil {
			return errTrans
		}

		// fmt.Println(txName, "create entry 1")
		result.FromEntry, errTrans = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if errTrans != nil {
			return errTrans
		}

		// fmt.Println(txName, "create entry 2")
		result.ToEntry, errTrans = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if errTrans != nil {
			return errTrans
		}

		// THIS IS THE OLD BLOCK
		// get account -> update its balance
		// fmt.Println(txName, "get account 1")
		// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }
		// // fmt.Println(txName, "update account 1")
		// result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.FromAccountID,
		// 	Balance: account1.Balance - arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }
		// // fmt.Println(txName, "get account 2")
		// account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		// if err != nil {
		// 	return err
		// }
		// // fmt.Println(txName, "update account 2")
		// result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.ToAccountID,
		// 	Balance: account2.Balance + arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		// THIS IS THE NEW BLOCK AND THE NEW WAY TO UPDATE ACCOUNT'S BALANCES
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, errTrans = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			return errTrans
		} else {
			result.ToAccount, result.FromAccount, errTrans = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			return errTrans
		}
		return nil
	})
	return result, err
}
func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return

}
