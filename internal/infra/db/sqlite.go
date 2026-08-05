package db

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func OpenSQLite(file string) *sql.DB {

	slog.Info("⏳ attempting to connect to the database.")

	db, err := sql.Open("sqlite3", file)
	if err != nil {

		slog.Error("❌ failed to connect to the database.", "error", err)
		os.Exit(1)

	}

	slog.Info("✅ successfully connected to the database.")
	slog.Info("⏳ attempting to ping the database.")

	if err := db.Ping(); err != nil {

		slog.Error("❌ ping to the database failed.", "error", err)
		os.Exit(1)

	}

	slog.Info("✅ successfully ping to the database.")
	slog.Info("⏳ enable the foreign_key")

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {

		slog.Error("🔛 enabling the foreign_keys failed.", "error", err)
		os.Exit(1)

	}

	slog.Info("✅ successfully enable the foreign_key")
	slog.Info("🚀 database is ready for use.")
	return db
}
