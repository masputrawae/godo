package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Open(file string) *sql.DB {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatal("error conneting to database: ", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("error ping to database: ", err)
	}

	if _, err = db.Exec("PRAGMA foreign_key = ON"); err != nil {
		log.Fatal("error enable foreign_key: ", err)
	}

	return db
}
