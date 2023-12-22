package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	DB *sql.DB
	*Queries
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		DB:      db,
		Queries: New(db),
	}
}

// Function to execute a transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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
