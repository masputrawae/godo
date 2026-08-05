package route

import (
	"godo/internal/handler"
	mw "godo/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Route struct {
	mux     *chi.Mux
	handler *handler.Handler
	mw      *mw.Middleware
}

func New(mux *chi.Mux, h *handler.Handler, mw *mw.Middleware) *Route {
	return &Route{
		mux:     mux,
		handler: h,
		mw:      mw,
	}
}

func (r *Route) Setup() {
	r.mux.Use(middleware.Logger, middleware.Recoverer)

	r.mux.HandleFunc("/register", r.handler.Register)
	r.mux.HandleFunc("/login", r.handler.Login)
	r.mux.HandleFunc("/logout", r.handler.Logout)

	r.mux.Route("/{username}", func(cr chi.Router) {
		cr.Use(r.mw.UserAuth)
		cr.Get("/", r.handler.User)
	})
}
