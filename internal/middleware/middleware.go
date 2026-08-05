package middleware

import (
	"context"
	"godo/internal/session"
	"net/http"
)

type Middleware struct {
	session *session.Session
}

func New(s *session.Session) *Middleware {
	return &Middleware{session: s}
}

func (m *Middleware) UserAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := m.session.Get(w, r)
		if session == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) CSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := m.session.Parse(r)
		if session == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		token := r.FormValue("csrf-token")
		if session.CSRFToken != token {
			w.WriteHeader(http.StatusForbidden)
		}

		next.ServeHTTP(w, r)
	})
}
