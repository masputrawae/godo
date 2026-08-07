package handler

import (
	"encoding/json"
	"godo/internal/model"
	"godo/internal/service"
	"net/http"
)

type Auth struct {
	auth    *service.Auth
	session *service.Session
}

// =====: Auth
func NewAuth(auth *service.Auth, session *service.Session) *Auth {
	return &Auth{auth: auth, session: session}
}

// =====: Login
func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	user, err := a.auth.Login(r.Context(), model.UserLoginPayload{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	})

	if err != nil {
		resErr(w, err, http.StatusUnauthorized)
	}

	if err := a.session.Set(r.Context(), w, user.ID); err != nil {
		resErr(w, err, http.StatusUnauthorized)
		return
	}

	jsonWrite(w, user)
}

// =====: Register
func (a *Auth) Register(w http.ResponseWriter, r *http.Request) {
	user, err := a.auth.Register(r.Context(), model.UserCreatePayload{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
	})

	if err != nil {
		resErr(w, err, http.StatusUnauthorized)
	}

	if err := a.session.Set(r.Context(), w, user.ID); err != nil {
		resErr(w, err, http.StatusUnauthorized)
		return
	}
	jsonWrite(w, user)
}

// =====: Response Error
func resErr(w http.ResponseWriter, err error, code int) {
	var res = struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    code,
		Message: err.Error(),
	}
	jsonWrite(w, res)
}

// =====: Json Writer
func jsonWrite(w http.ResponseWriter, data any) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
