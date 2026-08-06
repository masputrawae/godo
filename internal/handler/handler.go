package handler

import (
	"encoding/json"
	"godo/internal/cookie"
	"godo/internal/model"
	"godo/internal/service"
	"log"
	"net/http"
)

type Handler struct {
	cookie *cookie.Cookie
	user   *service.Service
}

func New(cookie *cookie.Cookie, user *service.Service) *Handler {
	return &Handler{
		cookie: cookie,
		user:   user,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	user, err := h.user.Create(r.Context(), model.UserCreatePayload{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionID, csrfToken, err := h.cookie.Set(r.Context(), w, user.ID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("SESSION ID: %s\nCSRF TOKEN: %s\n", sessionID, csrfToken)
	jsonW(w, user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	user, err := h.user.Login(r.Context(), model.UserLoginPayload{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionID, csrfToken, err := h.cookie.Set(r.Context(), w, user.ID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("SESSION ID: %s\nCSRF TOKEN: %s\n", sessionID, csrfToken)
	jsonW(w, user)
}

func (h *Handler) User(w http.ResponseWriter, r *http.Request) {
	session, ok := h.cookie.ParseCtx(r)
	if !ok {
		w.WriteHeader(http.StatusNonAuthoritativeInfo)
		return
	}

	user, err := h.user.GetByID(r.Context(), session.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNonAuthoritativeInfo)
		return
	}

	log.Println("WELCOME: ", user.Username)
	jsonW(w, user)
}

func jsonW(w http.ResponseWriter, data any) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
