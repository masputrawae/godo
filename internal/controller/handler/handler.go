package handler

import (
	"encoding/json"
	"errors"
	"godo/internal/controller/service"
	"godo/internal/infra/store"
	"godo/internal/model"
	"log/slog"
	"net/http"
)

type Handler struct {
	user    *service.User
	session *store.Session
}

func New(user *service.User, session *store.Session) *Handler {
	return &Handler{user: user, session: session}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		user, err := h.user.Register(r.Context(), model.UserCreatePayload{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
			Email:    r.FormValue("email"),
		})

		if err != nil {
			switch err {
			case service.ErrInvalidEmail, service.ErrUsernameToShort, service.ErrPasswordToShort:
				http.Error(w, err.Error(), http.StatusBadRequest)
			case service.ErrEmailAlreadyExists:
				http.Error(w, err.Error(), http.StatusConflict)
			default:
				slog.Error("INTERNAL SERVER ERROR", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		h.setCookie(w, user.ID)
		writeJson(w, user)
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		user, err := h.user.Login(r.Context(), model.UserLoginPayload{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		})

		if err != nil {
			if errors.Is(err, service.ErrUsernameOrPasswordNotMatch) {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}

			slog.Error("INTERNAL SERVER ERROR", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		h.setCookie(w, user.ID)
		writeJson(w, user)
		w.WriteHeader(http.StatusAccepted)
	}
}

// =====: HELPERS
func (h *Handler) setCookie(w http.ResponseWriter, userID int) {
	sessionID, err := h.session.GenID()
	if err != nil {
		slog.Error("INTERNAL SERVER ERROR", "error", err)
	}

	csrfToken, err := h.session.GenID()
	if err != nil {
		slog.Error("INTERNAL SERVER ERROR", "error", err)
	}

	h.session.SaveSession(userID, sessionID, csrfToken)
	h.session.SetCookie(w, sessionID)
}

func writeJson(w http.ResponseWriter, data any) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("INTERNAL SERVER ERROR", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
