package service

import (
	"context"
	"errors"
	"godo/internal/infra/password"
	"godo/internal/model"
	"godo/internal/repo"
)

var (
	ErrInvalidUsernameOrPassword = errors.New("Invalid username or password")
)

type User struct {
	rpUser *repo.User
}

func NewUser(rpUser *repo.User) *User {
	return &User{rpUser}
}

func (u *User) Register(ctx context.Context, p model.UserRequestRegister) (*model.User, error) {
	h, err := password.Hash(p.Password)
	if err != nil {
		return nil, err
	}

	p.Password = h
	user, err := u.rpUser.Create(ctx, p)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (u *User) Login(ctx context.Context, p model.UserRequestLogin) (*model.User, error) {
	user, err := u.rpUser.FindByUsername(ctx, p.Username)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, ErrInvalidUsernameOrPassword
		}
		return nil, err
	}

	if !password.Check(user.Password, p.Password) {
		return nil, ErrInvalidUsernameOrPassword
	}

	user.Password = ""
	return user, nil
}

func (u *User) Update(ctx context.Context, id int, p model.UserRequestUpdate) (*model.User, error) {
	if p.Password != nil {
		h, err := password.Hash(*p.Password)
		if err != nil {
			return nil, err
		}
		p.Password = &h
	}

	user, err := u.rpUser.Update(ctx, id, p)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	user, err := u.rpUser.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}
