package db

import "database/sql"

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
