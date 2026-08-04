package main

import (
	"fmt"
	"godo/internal/database"
	"godo/internal/repo"
)

func main() {
	db := database.ConnectSQLite("./app.db")
	defer db.Close()

	// ===== REPOSITORIES =====
	_ = repo.NewUser(db)
	_ = repo.NewTodo(db)
	_ = repo.NewStatus(db)
	_ = repo.NewPriority(db)

	fmt.Println("🐮 MOooo....")
}
