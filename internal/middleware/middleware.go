package middleware

import (
	"context"
	"godo/internal/model"
	"godo/internal/service"
	"log"
	"net/http"
	"time"
)

type Middleware struct {
	session *service.Session
}

func New(session *service.Session) *Middleware {
	return &Middleware{session: session}
}

func (m *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)
		log.Println(r.Method, "\t", time.Since(start), "\t", r.URL.Path)
	})
}

func (m *Middleware) SessionAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		session, err := m.session.Get(r.Context(), c.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if session.ExpiresAt.Before(time.Now()) {
			m.session.Delete(r.Context(), w, c.Value)
			http.Error(w, "session has expired", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *Middleware) CSRF(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, ok := r.Context().Value("session").(model.Session)
		if !ok {
			http.Error(w, "invalid session", http.StatusBadRequest)
			return
		}

		if s.CSRFToken != r.FormValue("csrf-token") {
			http.Error(w, "invalid csrf token", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	}
}
