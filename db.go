package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ctx context.Context
	db  *sql.DB
)

func DB() *sql.DB {
	database, _ := sql.Open("sqlite3", "./user.db")
	return database
}

func AfterTable() {
	tx, _ := DB().Begin()
	re, err := tx.Exec("ALTER TABLE users RENAME TO _users_old")
	if err != nil {
		tx.Rollback()
		fmt.Println(re, err)
	}
	re, err = tx.Exec(`CREATE TABLE users
		( 	id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name VARCHAR NOT NULL,
			last_name VARCHAR NOT NULL,
			email VARCHAR NOT NULL,
			gender VARCHAR NOT NULL,
			age INTEGER NOT NULL
	)`)
	if err != nil {
		tx.Rollback()
		fmt.Println(re, err)
	}

	re, err = tx.Exec(`INSERT INTO users (first_name,last_name,email,gender,age) SELECT first_name,last_name,email,gender,age FROM _users_old`)
	if err := tx.Commit(); err != nil {
		fmt.Println(re, err)
	}
	fmt.Println(re, err)
}
