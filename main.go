package main

import (
	"fmt"
	"godo/internal/database"
	"godo/internal/handler"
	"godo/internal/middleware"
	"godo/internal/repo"
	"godo/internal/route"
	"godo/internal/service"
	"godo/internal/session"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	db := database.ConnectSQLite("./app.db")
	defer db.Close()

	// ===== MUX =====
	mux := chi.NewMux()

	// ===== REPOSITORIES =====
	user := repo.NewUser(db)
	todo := repo.NewTodo(db)
	status := repo.NewStatus(db)
	priority := repo.NewPriority(db)

	// ===== SESSIONS =====
	session := session.New()

	// ===== MIDDLEWARES =====
	middleware := middleware.New(session)

	// ===== SERVICES =====
	todoService := service.NewTodo(todo, status, priority)
	userService := service.NewUser(user)

	// ===== HANDLERS =====
	handler := handler.New(userService, todoService, session)

	// ===== ROUTERS =====
	route.New(mux, handler, middleware).Setup()

	fmt.Println("🐮 MOooo....")
	http.ListenAndServe(":8080", mux)
}
