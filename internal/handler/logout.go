package handler

import (
	"log"
	"net/http"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(h.svSession.CookieName)
	if err == nil {
		if err := h.svSession.Remove(r.Context(), w, c.Value); err != nil {
			log.Println("error logout: ", err)
		}
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
