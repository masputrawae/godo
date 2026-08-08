package handler

import (
	"godo/internal/view"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	session, ok := h.svSession.ParseCtx(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := h.svUser.GetByID(r.Context(), session.UserID)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	isHXRequest, _ := strconv.ParseBool(r.Header.Get("HX-Request"))
	if err := view.Home(isHXRequest, user).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
