package service

import (
	"errors"
	"godo/internal/repo"
)

var (
	ErrEmailAlreadyExists         = errors.New("email already exists")
	ErrUsernameAlreadyExists      = errors.New("email already exists")
	ErrInvalidEmail               = errors.New("invalid email")
	ErrUsernameToShort            = errors.New("username to short. min 3 character")
	ErrPasswordToShort            = errors.New("password to short. min 8 character")
	ErrUsernameOrPasswordNotMatch = errors.New("username or password not match")
)

type User struct {
	user *repo.User
}
