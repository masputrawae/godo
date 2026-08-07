package middleware

import (
	"godo/internal/service"
	"log"
	"net/http"
	"time"
)

type Middleware struct {
	svSession *service.Session
}

func New(svSession *service.Session) *Middleware {
	return &Middleware{svSession}
}

func (m *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s\t%s\t%v\n", r.Method, r.URL.Path, time.Since(start))
	})
}
