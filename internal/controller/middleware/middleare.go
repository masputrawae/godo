package middleware

import (
	"context"
	"fmt"
	"godo/internal/controller/cookie"
	"log"
	"net/http"
	"time"
)

type Middleware struct {
	session *cookie.Session
}

func New(session *cookie.Session) *Middleware {
	return &Middleware{session: session}
}

func (m *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path, time.Since(start))
		if r.Method == "POST" {
			fmt.Println("DEBUG FORM")
			fmt.Println("Email:", r.FormValue("username"))
			fmt.Println("Username:", r.FormValue("username"))
			fmt.Println("Password:", r.FormValue("username"))
		}
	})
}

func (m *Middleware) SessionAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, ok := m.session.Get(r)
		if !ok {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "session", *session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) CSRFAuth() {

}
