package repo

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"alukart32.com/bank/config"
	"alukart32.com/bank/pkg/pqx"
)

var (
	testDB     *sql.DB
	testDBConf = &config.DB{
		URI:    "postgres://postgres:postgres@localhost:5432/bank-test?sslmode=disable",
		Driver: "postgres",
	}
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = pqx.New(testDBConf)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	os.Exit(m.Run())
}
