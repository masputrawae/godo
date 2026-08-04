package main

import (
	"fmt"
	"godo/internal/database"
	"godo/internal/repo"
	"godo/internal/service"
)

func main() {
	db := database.ConnectSQLite("./app.db")
	defer db.Close()

	// ===== REPOSITORIES =====
	user := repo.NewUser(db)
	todo := repo.NewTodo(db)
	status := repo.NewStatus(db)
	priority := repo.NewPriority(db)

	// ===== SERVICES =====
	_ = service.NewTodo(todo, status, priority)
	_ = service.NewUser(user)

	fmt.Println("🐮 MOooo....")
}
