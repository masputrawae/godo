package handler

import (
	"encoding/json"
	"godo/internal/model"
	"godo/internal/service"
	"godo/internal/session"
	"log"
	"net/http"
)

type Handler struct {
	user    *service.User
	todo    *service.Todo
	session *session.Session
}

func New(u *service.User, t *service.Todo, s *session.Session) *Handler {
	return &Handler{user: u, todo: t, session: s}
}

func (s *Handler) Root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome"))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Login Page"))

	case http.MethodPost:
		user, err := h.user.Login(r.Context(), model.UserLoginPayload{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		})

		if err != nil {
			switch err {
			case service.ErrUsernameOrPasswordNotMatch, service.ErrUserAccountNotFound:
				http.Error(w, service.ErrUsernameOrPasswordNotMatch.Error(), http.StatusUnauthorized)
			default:
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("INTERNAL SERVER ERROR: ", err)
			}
			return
		}

		if err := h.session.Set(w, user.ID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("INTERNAL SERVER ERROR: ", err)
		}

		if err := json.NewEncoder(w).Encode(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("INTERNAL SERVER ERROR: ", err)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Register Page"))
	case http.MethodPost:
		user, err := h.user.Register(r.Context(), model.UserCreatePayload{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
			Email:    r.FormValue("email"),
		})

		if err != nil {
			switch err {
			case
				service.ErrEmailAlreadyUse,
				service.ErrEmailInvalid,
				service.ErrUsernameAlreadyUse,
				service.ErrUsernameToShort,
				service.ErrPasswordToShort:
				http.Error(w, err.Error(), http.StatusBadRequest)
			default:
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("INTERNAL SERVER ERROR: ", err)
			}
			return
		}

		if err := h.session.Set(w, user.ID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("INTERNAL SERVER ERROR: ", err)
		}

		if err := json.NewEncoder(w).Encode(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("INTERNAL SERVER ERROR: ", err)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) User(w http.ResponseWriter, r *http.Request) {
	session := h.session.Parse(r)
	if session == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	id := session.UserID
	user, err := h.user.GetByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]any{"INFO": user}); err != nil {
		log.Println(err)
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	id, err := r.Cookie("session")
	if err != nil {
		return
	}

	h.session.Delete(w, id.Value)
}
