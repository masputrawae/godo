package service

import (
	"context"
	"errors"
	"godo/internal/infra/password"
	"godo/internal/model"
	"godo/internal/repo"
	"net/mail"
)

var (
	ErrUsernameToShort            = errors.New("username to short. min 3 character")
	ErrPasswordToShort            = errors.New("password to short. min 8 character")
	ErrUsernameOrPasswordNotMatch = errors.New("username or password not match")
)

type Auth struct {
	user *repo.User
}

// =====: Auth
func NewAuth(user *repo.User) *Auth {
	return &Auth{user}
}

// =====: Register
func (a *Auth) Register(ctx context.Context, p model.UserCreatePayload) (model.User, error) {
	if len(p.Username) < 3 {
		return model.User{}, ErrUsernameToShort
	}
	if len(p.Password) < 8 {
		return model.User{}, ErrPasswordToShort
	}
	if _, err := mail.ParseAddress(p.Email); err != nil {
		return model.User{}, err
	}

	h, err := password.Hash(p.Password)
	if err != nil {
		return model.User{}, err
	}

	p.Password = h
	id, err := a.user.Create(ctx, p)
	if err != nil {
		return model.User{}, err
	}

	user, err := a.user.FindByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	user.Password = ""
	return user, nil
}

// =====: Login
func (a *Auth) Login(ctx context.Context, p model.UserLoginPayload) (model.User, error) {
	u, err := a.user.FindByUsername(ctx, p.Username)
	if err != nil {
		return model.User{}, ErrUsernameOrPasswordNotMatch
	}

	if !password.Check(u.Password, p.Password) {
		return model.User{}, ErrUsernameOrPasswordNotMatch
	}

	u.Password = ""
	return u, nil
}

// =====: Helpers
