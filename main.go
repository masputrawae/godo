package main

import (
	"godo/internal/controller/service"
	"godo/internal/infra/db"
	_ "godo/internal/infra/logger"
	"godo/internal/repo"
)

func main() {
	db := db.OpenSQLite("app.db")
	defer db.Close()

	// ===== REPOSITORIES =====
	repo := repo.New(db)
	repoUser := repo.NewUser()
	/*
		repoTodo := repo.NewTodo()
		repoStatus := repo.NewStatus()
		repoPriority := repo.NewPriority()
	*/

	serviceUser := service.NewUser(repoUser)
}
