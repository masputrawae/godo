package handler

import (
	"encoding/json"
	"godo/internal/model"
	"godo/internal/service"
	"log"
	"net/http"
)

type Handler struct {
	user *service.User
	todo *service.Todo
}

func New(u *service.User, t *service.Todo) *Handler {
	return &Handler{user: u, todo: t}
}

func (p *Handler) Root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome"))
}

func (p *Handler) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Login Page"))

	case http.MethodPost:
		user, err := p.user.Login(r.Context(), model.UserLoginPayload{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		})

		if err != nil {
			switch err {
			case service.ErrUsernameOrPasswordNotMatch, service.ErrUserAccountNotFound:
				http.Error(w, service.ErrUsernameOrPasswordNotMatch.Error(), http.StatusUnauthorized)
			default:
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
			}
			return
		}

		if err := json.NewEncoder(w).Encode(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (p *Handler) Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Register Page"))
	case http.MethodPost:
		user, err := p.user.Register(r.Context(), model.UserCreatePayload{
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
				log.Println(err)
			}
			return
		}

		if err := json.NewEncoder(w).Encode(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
