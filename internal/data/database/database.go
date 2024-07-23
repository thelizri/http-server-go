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

	if dburl == "" {
		dburl = ":memory:"
	}

	db, err = sql.Open("sqlite3", dburl)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}

	initTables()
}

func initTables() {
	db.Exec(`CREATE TABLE IF NOT EXISTS user (
            id INTEGER NOT NULL PRIMARY KEY ASC, 
            username TEXT NOT NULL UNIQUE, 
            password TEXT NOT NULL CHECK(length(password) > 5))`)
}
