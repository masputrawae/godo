package middleware

import (
	"context"
	"godo/internal/cookie"
	"log"
	"net/http"
	"time"
)

type Middleware struct {
	cookie *cookie.Cookie
}

func New(cookie *cookie.Cookie) *Middleware {
	return &Middleware{cookie: cookie}
}

func (m *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)
		log.Println(r.Method, "\t", r.URL.Path, "\t", time.Since(start))
	})
}

func (m *Middleware) SessionAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(m.cookie.Name)
		if err != nil {
			w.WriteHeader(http.StatusNonAuthoritativeInfo)
			return
		}

		session, err := m.cookie.Get(r.Context(), w, cookie.Value)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNonAuthoritativeInfo)
			return
		}

		ctx := context.WithValue(r.Context(), m.cookie.Name, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
