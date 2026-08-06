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
	ErrUsernameToShort            = errors.New("username to short")
	ErrPasswordToShort            = errors.New("password to short")
	ErrEmailInvalid               = errors.New("email invalid")
	ErrUsernameAlreadyUse         = errors.New("username is already use")
	ErrEmailAlreadyUse            = errors.New("email is already use")
	ErrUsernameOrPasswordNotMatch = errors.New("username or password not match")
	ErrUserNotFound               = errors.New("user not found")
)

type User struct {
	user *repo.User
}

func NewUser(user *repo.User) *User {
	return &User{user: user}
}

// =====: CREATE USER
func (u *User) Create(ctx context.Context, p model.UserCreatePayload) (*int, error) {
	if len(p.Username) < 3 {
		return nil, ErrUsernameToShort
	}
	if len(p.Password) < 8 {
		return nil, ErrPasswordToShort
	}
	if _, err := mail.ParseAddress(p.Email); err != nil {
		return nil, ErrEmailInvalid
	}

	hash, err := password.Hash(p.Password)
	if err != nil {
		return nil, err
	}

	p.Password = *hash
	id, err := u.user.Create(ctx, p)
	if err != nil {
		switch u.user.TranslateError(err) {
		case repo.ErrUniqueEmail:
			return nil, ErrEmailAlreadyUse
		case repo.ErrUniqueUsername:
			return nil, ErrUsernameAlreadyUse
		}

		return nil, err
	}

	return id, nil
}

// =====: LOGIN USER
func (u *User) Login(ctx context.Context, p model.UserLoginPayload) (*model.User, error) {
	user, err := u.user.FindByUsername(ctx, p.Username)
	if err != nil {
		switch u.user.TranslateError(err) {
		case repo.ErrUserNotFound:
			return nil, ErrUsernameOrPasswordNotMatch
		}

		return nil, err
	}

	if !password.Check(user.Password, p.Password) {
		return nil, ErrUsernameOrPasswordNotMatch
	}

	return user, nil
}

// =====: GET USER BY ID
func (u *User) GetByID(ctx context.Context, id int) (*model.User, error) {
	user, err := u.user.FindByID(ctx, id)
	if err != nil {
		switch u.user.TranslateError(err) {
		case repo.ErrUserNotFound:
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user.Password = ""
	return user, nil
}
