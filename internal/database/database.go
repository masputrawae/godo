package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectSQLite(file string) *sql.DB {
	log.Println("⏳ create connection to sqlite.")

	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatal("❌ connection error.\n", err)
	}

	log.Println("🔵 test ping to sqlite.")
	if err := db.Ping(); err != nil {
		log.Fatal("🔴 failed ping\n", err)
	}

	log.Println("🔛 enabling PRAGMA foreign_keys.")
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		log.Fatal(err)
	}

	log.Println("✅ database is connected. it is now ready to use.")
	return db
}
