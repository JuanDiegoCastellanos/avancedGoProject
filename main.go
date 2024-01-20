package main

import (
	"avancedGo/api"
	db "avancedGo/db/sqlc"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var (
	dbDriver = "postgres"
	dbSource = "postgresql://root:manolo221212@localhost:5433/simple_posts?sslmode=disable"
	address  = "0.0.0.0:8080"
)

func main() {
	testDB, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to database... ", err)
	}
	store := db.NewStore(testDB)
	server := api.NewServer(store)

	err = server.Start(address)
	if err != nil {
		log.Fatal("cannot start the server... ", err)
	}
}
