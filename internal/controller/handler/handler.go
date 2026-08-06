package handler

import (
	"encoding/json"
	"godo/internal/controller/cookie"
	"godo/internal/controller/service"
	"godo/internal/model"
	"log/slog"
	"net/http"
)

type Handler struct {
	user    *service.User
	session *cookie.Session
}

func New(user *service.User, session *cookie.Session) *Handler {
	return &Handler{user: user, session: session}
}

// =====: REGISTER
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	id, err := h.user.Create(r.Context(), model.UserCreatePayload{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
	})

	if err != nil {
		switch err {
		case
			service.ErrEmailAlreadyUse,
			service.ErrUsernameAlreadyUse:
			http.Error(w, err.Error(), http.StatusNonAuthoritativeInfo)
		case
			service.ErrPasswordToShort,
			service.ErrUsernameToShort,
			service.ErrPasswordToShort:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			slog.Error("Register", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	user, err := h.user.GetByID(r.Context(), *id)
	if err != nil {
		slog.Error("Register", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionID, err := h.session.GenID()
	if err != nil {
		slog.Error("Generate Session ID", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	csrfToken, err := h.session.GenID()
	if err != nil {
		slog.Error("Generate CSRF Token", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.session.Set(w, user.ID, *sessionID, *csrfToken)
	writeJson(w, user)
	w.WriteHeader(http.StatusAccepted)
}

// =====: LOGIN
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	user, err := h.user.Login(r.Context(), model.UserLoginPayload{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	})

	if err != nil {
		switch err {
		case service.ErrUsernameOrPasswordNotMatch:
			w.WriteHeader(http.StatusForbidden)
			return
		}

		slog.Error("Login", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionID, err := h.session.GenID()
	if err != nil {
		slog.Error("Generate Session ID", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	csrfToken, err := h.session.GenID()
	if err != nil {
		slog.Error("Generate CSRF Token", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.session.Set(w, user.ID, *sessionID, *csrfToken)
	writeJson(w, user)
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) User(w http.ResponseWriter, r *http.Request) {
	session, ok := h.session.Parse(r)
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	user, err := h.user.GetByID(r.Context(), session.UserID)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	writeJson(w, user)
}

// =====: WRITE JSON
func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Write Json", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
