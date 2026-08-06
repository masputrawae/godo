package main

import (
	"godo/internal/cookie"
	"godo/internal/handler"
	"godo/internal/infra/db"
	"godo/internal/middleware"
	"godo/internal/repo"
	"godo/internal/service"
	"log"
	"net/http"
)

func main() {
	// =====: Store
	db := db.OpenSQLite("app.db")
	defer db.Close()

	// =====: Repositories
	repo := repo.New(db)
	repoSession := repo.NewSession()
	repoUser := repo.NewUser()

	repoSession.CleanUp()

	// =====: Services
	serviceUser := service.New(repoUser)

	// =====: Mux, Cookie Middlewares, Handlers
	mux := http.NewServeMux()
	cookie := cookie.New(repoSession)
	mw := middleware.New(cookie)
	handler := handler.New(cookie, serviceUser)

	// =====: Route
	mux.Handle("/", mw.SessionAuth(http.HandlerFunc(handler.User)))
	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/register", handler.Register)

	log.Fatal(http.ListenAndServe(":8080", mw.Logger(mux)))
}
