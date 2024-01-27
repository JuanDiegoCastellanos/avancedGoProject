package main

import (
	"github.com/JuanDiegoCastellanos/advancedGoProject/api"
	db "github.com/JuanDiegoCastellanos/advancedGoProject/db/sqlc"
	"github.com/JuanDiegoCastellanos/advancedGoProject/util"

	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot read the app config file", err)
	}
	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database... ", err)
	}
	store := db.NewStore(testDB)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start the server... ", err)
	}
}
