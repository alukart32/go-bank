package pqx

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"alukart32.com/bank/config"
	_ "github.com/lib/pq"
)

// db connection pool
var (
	db   *sql.DB
	once sync.Once
)

// New returns a new sql.DB if it was not installed earlier,
// otherwise it returns an existing instance. An error occurs
// if the database connection wasn't initialized.
func New(cfg *config.DB) (_ *sql.DB, err error) {
	once.Do(func() {
		log.Printf("init the new db connection pool...")
		db, err = sql.Open(cfg.Driver, cfg.URI)
		if err != nil {
			err = fmt.Errorf("%s database open error %w", cfg.Driver, err)
		}

		if err = db.Ping(); err == nil {
			log.Printf("db connection pool was created...")
		}
	})
	return db, err
}

// Close closes the database connection pool.
func Close() {
	if db != nil {
		db.Close()
	}
}
