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
	testDB   *sql.DB
	testConf *config.Config

	// &config.DB{
	// 	URI:    "postgres://postgres:postgres@localhost:5432/bank-test?sslmode=disable",
	// 	Driver: "postgres",
	// }
)

func TestMain(m *testing.M) {
	var err error
	testConf, err = config.New("test")
	if err != nil {
		log.Fatal("cannot get config: ", err)
	}

	testDB, err = pqx.New(&testConf.DB)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	os.Exit(m.Run())
}
