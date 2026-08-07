package handler

import "net/http"

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	session, ok := h.svSession.ParseCtx(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	_, err := h.svUser.GetByID(r.Context(), session.UserID)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}
