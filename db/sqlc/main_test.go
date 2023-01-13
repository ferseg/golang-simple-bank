package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:sa123@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
  var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to the datbase", err)
	}
	testQueries = New(testDB)

	m.Run()
}
