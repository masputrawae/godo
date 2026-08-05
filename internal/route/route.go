package route

import (
	"godo/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Route struct {
	mux     *chi.Mux
	handler *handler.Handler
}

func New(mux *chi.Mux, h *handler.Handler) *Route {
	return &Route{
		mux:     mux,
		handler: h,
	}
}

func (r *Route) Setup() {
	r.mux.Use(middleware.Logger, middleware.Recoverer)

	r.mux.HandleFunc("/register", r.handler.Register)
	r.mux.HandleFunc("/login", r.handler.Login)
}
