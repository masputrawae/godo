package main

import (
	"context"
	"godo/internal/handler"
	"godo/internal/infra/db"
	"godo/internal/middleware"
	"godo/internal/repo"
	"godo/internal/service"
	"log"
	"net/http"
	"time"
)

func autoCleanSession(session *repo.Session, ticker *time.Ticker) {
	for range ticker.C {
		if err := session.DeleteByExpires(context.Background()); err != nil {
			log.Println(err)
		}
	}
}

func main() {
	ticker := time.NewTicker(30 * time.Second)

	// =====: Storage
	db := db.Open("app.db")
	defer db.Close()

	// =====: Repositories
	rp := repo.New(db)
	rpUser := repo.NewUser(rp)
	rpSession := repo.NewSession(rp)

	// =====: auto clean session
	go autoCleanSession(rpSession, ticker)

	// =====: Services
	svAuth := service.NewAuth(rpUser)
	svSession := service.NewSession(rpSession)

	// =====: Middlewares
	mw := middleware.New(svSession)

	// =====: Handlers
	hlAuth := handler.NewAuth(svAuth, svSession)

	// =====: Mux
	mux := http.NewServeMux()

	mux.HandleFunc("/login", hlAuth.Login)
	mux.HandleFunc("/register", hlAuth.Register)

	log.Fatal(http.ListenAndServe(":8080", mw.Logger(mux)))
}
