package main

import (
	"context"
	"godo/internal/handler"
	"godo/internal/infra/database"
	"godo/internal/middleware"
	"godo/internal/repo"
	"godo/internal/service"
	"log"
	"net/http"
	"time"
)

func main() {
	db := database.Open("app.db")

	// ==========: Repositories
	rpUser := repo.NewUser(db)
	rpTodo := repo.NewTodo(db)
	rpStatus := repo.NewStatus(db)
	rpPriority := repo.NewPriority(db)
	rpSession := repo.NewSession(db)

	// ==========: Services
	svUser := service.NewUser(rpUser)
	svSession := service.NewSession(rpSession, 1*time.Minute, "session")
	svTodo := service.NewTodo(rpTodo, rpPriority, rpStatus)

	// ==========: Auto Clean (sessions)
	go svSession.AutoClean(context.Background(), nil)

	// ==========: Handler
	handler := handler.New(svUser, svSession, svTodo)

	// ==========: Middleware
	mw := middleware.New(svSession)

	// ==========: Mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", mw.SessionAuth(handler.Root))
	mux.HandleFunc("/logout", mw.SessionAuth(handler.Logout))

	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/register", handler.Register)

	log.Println("server running on :8080")
	http.ListenAndServe(":8080", mw.Logger(mux))
}
