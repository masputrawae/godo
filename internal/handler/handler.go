package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"godo/internal/model"
	"godo/internal/repo"
	"godo/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	validator *validator.Validate
	svUser    *service.User
}

func New(validator *validator.Validate, svUser *service.User) *Handler {
	return &Handler{validator: validator, svUser: svUser}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	req := model.UserRequestRegister{
		Email:    *formString(r, "email"),
		Username: *formString(r, "username"),
		Password: *formString(r, "password"),
	}

	if err := h.validator.Struct(req); err != nil {
		if ve, ok := errors.AsType[validator.ValidationErrors](err); ok {
			for _, e := range ve {
				http.Error(w, fmt.Sprintf("%s %s %s", e.Field(), e.Tag(), e.Value()), http.StatusBadRequest)
				return
			}
		}

		log.Println("error internal validator:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user, err := h.svUser.Register(r.Context(), req)
	if err != nil {
		switch err {
		case repo.ErrUserAlreadyUsed, repo.ErrEmailAlreadyUsed, repo.ErrUsernameAlreadyUsed:
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		log.Println("error internal server register:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// TODO: 2026-08-08 - 2026-08-10
	// implementasi JWT nanti aja. kalau ini udah testing dan jalan
	// ....

	// sementara cek pakai curl (belum buat ui), nanti render pakai templ + htmx
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	req := model.UserRequestLogin{
		Username: *formString(r, "username"),
		Password: *formString(r, "password"),
	}

	if err := h.validator.Struct(req); err != nil {
		if ve, ok := errors.AsType[validator.ValidationErrors](err); ok {
			for _, e := range ve {
				http.Error(w, fmt.Sprintf("%s %s %s", e.Field(), e.Tag(), e.Value()), http.StatusBadRequest)
				return
			}
		}

		log.Println("error internal validator:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user, err := h.svUser.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUsernameOrPassword) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		log.Println("error internal server login:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	req := model.UserRequestUpdate{
		Email:    formString(r, "email"),
		Username: formString(r, "username"),
		Password: formString(r, "password"),
	}

	if err := h.validator.Struct(req); err != nil {
		if ve, ok := errors.AsType[validator.ValidationErrors](err); ok {
			for _, e := range ve {
				http.Error(w, fmt.Sprintf("%s %s %s", e.Field(), e.Tag(), e.Value()), http.StatusBadRequest)
				return
			}
		}

		log.Println("error internal validator:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user, err := h.svUser.Update(r.Context(), id, req)
	if err != nil {
		switch err {
		case repo.ErrUserAlreadyUsed, repo.ErrEmailAlreadyUsed, repo.ErrUsernameAlreadyUsed:
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		log.Println("error internal server update:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// convert form value
func formString(r *http.Request, key string) *string {
	value, ok := r.Form[key]
	if !ok {
		return nil
	}

	return &value[0]
}
