package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dburl                = os.Getenv("DB_URL")
	db                   *sql.DB
	dbRepositoryInstance *dbRepository
)

func init() {
	var err error

	db, err = sql.Open("sqlite3", dburl)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}
}
