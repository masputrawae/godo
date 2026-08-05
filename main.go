package main

import (
	"godo/internal/infra/db"
	_ "godo/internal/infra/logger"
)

func main() {
	db := db.OpenSQLite("app.db")
	defer db.Close()
}
