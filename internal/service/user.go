package service

import (
	"context"
	"database/sql"
	"errors"
	"godo/internal/infra/password"
	"godo/internal/model"
	"godo/internal/repo"
	"net/mail"
	"strings"
)

var (
	ErrUsernameToShort             = errors.New("Username to short")
	ErrPasswordToShort             = errors.New("Password to short")
	ErrUsernameAlreadyUse          = errors.New("Username already in use")
	ErrEmailAlreadyUse             = errors.New("Email already in use")
	ErrIncorrectUsernameOrPassword = errors.New("Incorrect username or password")
)

type User struct {
	rpUser *repo.User
}

func NewUser(rpUser *repo.User) *User {
	return &User{rpUser}
}

func (u *User) Register(ctx context.Context, p model.UserPayloadCreate) (*model.User, error) {
	if len(p.Username) < 3 {
		return nil, ErrUsernameToShort
	}

	if len(p.Password) < 8 {
		return nil, ErrUsernameToShort
	}

	if _, err := mail.ParseAddress(p.Email); err != nil {
		return nil, err
	}

	h, err := password.Hash(p.Password)
	if err != nil {
		return nil, err
	}
	p.Password = h

	err = u.rpUser.Create(ctx, p)
	if err != nil {
		se := strings.ToLower(err.Error())
		switch {
		case strings.Contains(se, "unique"), strings.Contains(se, "users.username"):
			return nil, ErrUsernameAlreadyUse
		case strings.Contains(se, "unique"), strings.Contains(se, "users.email"):
			return nil, ErrEmailAlreadyUse
		default:
			return nil, err
		}
	}

	user, err := u.rpUser.FindByUsername(ctx, p.Username)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return &user, nil
}

func (u *User) Login(ctx context.Context, p model.UserPayloadLogin) (*model.User, error) {
	user, err := u.rpUser.FindByUsername(ctx, p.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrIncorrectUsernameOrPassword
		}
		return nil, err
	}

	if !password.Check(user.Password, p.Password) {
		return nil, ErrIncorrectUsernameOrPassword
	}

	user.Password = ""
	return &user, nil
}

func (u *User) GetByID(ctx context.Context, id int) (model.User, error) {
	user, err := u.rpUser.FindByID(ctx, id)
	if err != nil {
		return user, err
	}

	user.Password = ""
	return user, nil
}
