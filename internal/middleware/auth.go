package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

func (m *Middleware) SessionAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(m.svSession.CookieName)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			log.Println("Session auth error:", err)
			return
		}

		session, err := m.svSession.Get(r.Context(), c.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			log.Println("Session auth error:", err)
			return
		}

		ctx := context.WithValue(r.Context(), m.svSession.CookieName, session)
		next(w, r.WithContext(ctx))
	}
}

func (m *Middleware) CSRFAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			next(w, r)
		}

		session, ok := m.svSession.ParseCtx(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			log.Println("CSRF auth error: parse ctx")
			return
		}

		token := r.FormValue("csrf-token")
		if session.CSRFToken != token {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if session.ExpiresAt.Before(time.Now()) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}
