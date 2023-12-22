package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
	dbDriver    = "postgres"
	dbSource    = "postgresql://root:manolo221212@localhost:5433/simple_posts?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
