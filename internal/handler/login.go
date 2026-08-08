package handler

import (
	"godo/internal/model"
	"godo/internal/view"
	"net/http"
	"strconv"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(h.svSession.CookieName)
	if err == nil {
		_, err := h.svSession.Get(r.Context(), c.Value)
		if err == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}

	switch r.Method {
	case http.MethodPost:

		user, err := h.svUser.Login(r.Context(), model.UserPayloadLogin{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err := h.svSession.Set(r.Context(), w, user.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Replace-Url", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	case http.MethodGet:
		isHXRequest, _ := strconv.ParseBool(r.Header.Get("HX-Request"))
		view.Login(isHXRequest).Render(r.Context(), w)
	}
}
