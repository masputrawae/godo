package handler

import (
	"godo/internal/model"
	"net/http"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(h.svSession.CookieName)
	if err == nil {
		_, err := h.svSession.Get(r.Context(), c.Value)
		if err == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}

	switch r.Method {
	case http.MethodPost:
		user, err := h.svUser.Register(r.Context(), model.UserPayloadCreate{
			Email:    r.FormValue("email"),
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.svSession.Set(r.Context(), w, user.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodGet:
	}
}
