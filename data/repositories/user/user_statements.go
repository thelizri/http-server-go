package userrepository

import (
	"database/sql"
	"http-server/data/database"
	"log"
)

var (
	dbRepository    database.DbRepository
	createUserStmt  *sql.Stmt
	getUserByIdStmt *sql.Stmt
)

func init() {
	dbRepository = database.NewDbRepository()
	prepareStatements()
}

func prepareStatements() {
	prepareCreateUserStmt()
	prepareGetUserByIdStmt()
}

func prepareCreateUserStmt() {
	query := "INSERT INTO user (username, password) VALUES (?, ?)"

	if stmt, err := dbRepository.Prepare(query); err != nil {
		log.Fatal("Could not prepare Create User statement: ", err)
	} else {
		createUserStmt = stmt
	}
}

func prepareGetUserByIdStmt() {
	query := "SELECT * FROM user WHERE ID = ?"

	if stmt, err := dbRepository.Prepare(query); err != nil {
		log.Fatal("Could not prepare Get User By Id statement: ", err)
	} else {
		getUserByIdStmt = stmt
	}
}
