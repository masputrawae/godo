package service

import (
	"context"
	"errors"
	"godo/internal/infra/password"
	"godo/internal/model"
	"godo/internal/repo"
)

// =====: USER LOGIN
func (u *User) Login(ctx context.Context, p model.UserLoginPayload) (*model.User, error) {
	user, err := u.user.FindByUsername(ctx, p.Username)
	if err != nil {
		if errors.Is(err, repo.ErrDataNotFound) {
			return nil, ErrUsernameOrPasswordNotMatch
		}

		return nil, err
	}

	if !password.Check(user.Password, p.Password) {
		return nil, ErrUsernameOrPasswordNotMatch
	}

	return user, err
}
