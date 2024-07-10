package database

import (
	"database/sql"
	"log"
)

var (
	createUserStmt  *sql.Stmt
	getUserByIdStmt *sql.Stmt
)

func prepareUserStatements() {
	prepareCreateUserStmt()
	prepareGetUserByIdStmt()
}

func prepareCreateUserStmt() {
	query := "INSERT INTO user (username, password) VALUES (?, ?)"

	if stmt, err := db.Prepare(query); err != nil {
		log.Fatal("Could not prepare Create User statement")
	} else {
		createUserStmt = stmt
	}
}

func prepareGetUserByIdStmt() {
	query := "SELECT * FROM user WHERE ID = ?"

	if stmt, err := db.Prepare(query); err != nil {
		log.Fatal("Could not prepare Get User By Id statement")
	} else {
		getUserByIdStmt = stmt
	}
}
