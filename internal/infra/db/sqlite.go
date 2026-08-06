package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func OpenSQLite(file string) *sql.DB {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatal(err)
	}

	return db
}
