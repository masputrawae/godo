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
	ErrUsernameToShort             = errors.New("username to short. min 3 character")
	ErrPasswordToShort             = errors.New("password to short. min 8 character")
	ErrEmailInvalid                = errors.New("email invalid")
	ErrUsernameAndPasswordNotMatch = errors.New("username and password not match")
	ErrUsernameAlreadyExist        = errors.New("username already exists")
	ErrEmailAlreadyExist           = errors.New("email already exists")
	ErrDataUserNotFound            = errors.New("data user not found")
)

type Service struct {
	user *repo.User
}

// =====: New
func New(user *repo.User) *Service {
	return &Service{user: user}
}

// =====: Create
func (s *Service) Create(ctx context.Context, p model.UserCreatePayload) (model.User, error) {
	if len(p.Username) < 3 {
		return model.User{}, ErrUsernameToShort
	}
	if len(p.Password) < 8 {
		return model.User{}, ErrPasswordToShort
	}
	if _, err := mail.ParseAddress(p.Email); err != nil {
		return model.User{}, ErrEmailInvalid
	}

	h, err := password.Hash(p.Password)
	if err != nil {
		return model.User{}, err
	}

	p.Password = h
	id, err := s.user.Create(ctx, p)
	if err != nil {
		switch err {
		case repo.ErrUsernameAlreadyExist:
			return model.User{}, ErrUsernameAlreadyExist
		case repo.ErrEmailAlreadyExist:
			return model.User{}, ErrEmailAlreadyExist
		default:
			return model.User{}, err
		}
	}

	return s.GetByID(ctx, id)
}

// =====: Login
func (s *Service) Login(ctx context.Context, p model.UserLoginPayload) (model.User, error) {
	if p.Username == "" && p.Password == "" {
		return model.User{}, ErrUsernameAndPasswordNotMatch
	}

	user, err := s.user.FindByUsername(ctx, p.Username)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrUserNotFound):
			return user, ErrUsernameAndPasswordNotMatch
		default:
			return user, err
		}
	}

	if !password.Check(user.Password, p.Password) {
		return model.User{}, ErrUsernameAndPasswordNotMatch
	}

	user.Password = ""
	return user, nil
}

// =====: Get
func (s *Service) GetByID(ctx context.Context, userID int) (model.User, error) {
	user, err := s.user.FindByID(ctx, userID)
	if err != nil {
		switch err {
		case repo.ErrUserNotFound:
			return user, ErrUsernameAndPasswordNotMatch
		default:
			return user, err
		}
	}

	user.Password = ""
	return user, nil
}
