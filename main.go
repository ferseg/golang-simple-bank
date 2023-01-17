package main

import (
	"database/sql"
	"log"

	"github.com/ferseg/golang-simple-bank/api"
	db "github.com/ferseg/golang-simple-bank/db/sqlc"
	"github.com/ferseg/golang-simple-bank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Something has happend", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to the db: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Could not create the server", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start the server", err)

	}

}
