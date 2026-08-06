package main

import (
	"godo/internal/controller/cookie"
	"godo/internal/controller/handler"
	"godo/internal/controller/middleware"
	"godo/internal/controller/service"
	"godo/internal/infra/db"
	"godo/internal/infra/store"
	"godo/internal/repo"
	"log"
	"net/http"
)

func main() {
	// ===== STORES =====
	db := db.OpenSQLite("app.db")
	defer db.Close()
	storeSession := store.NewSession()

	// ===== REPOSITORIES =====
	repo := repo.New(db)
	repoUser := repo.NewUser()
	/*
		repoTodo := repo.NewTodo()
		repoStatus := repo.NewStatus()
		repoPriority := repo.NewPriority()
	*/

	// ===== SERVICES =====
	serviceUser := service.NewUser(repoUser)

	// ===== COOKIES =====
	cookieSession := cookie.NewSession(storeSession)

	// ===== MUX, HANDLERS, MIDDLEWARES =====
	mux := http.NewServeMux()
	handler := handler.New(serviceUser, cookieSession)
	middleware := middleware.New(cookieSession)

	mux.HandleFunc("/", middleware.SessionAuth(handler.User))
	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/register", handler.Register)

	log.Fatal(http.ListenAndServe(":8080", middleware.Logger(mux)))
}
