package main

import (
	"godo/internal/infra/database"
)

func main() {
	// validator
	// validate := validator.New(validator.WithRequiredStructEnabled())

	// database (sqlite3)
	db := database.Open("app.db")
	defer db.Close()

	// repositories
	// rpUser := repo.NewUser(db)

	// services
	// svUser := service.NewUser(rpUser)

}
