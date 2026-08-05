package service

import (
	"context"
	"godo/internal/infra/password"
	"godo/internal/model"
	"net/mail"
	"unicode/utf8"
)

// =====: USER REGISTER
func (u *User) Register(ctx context.Context, p model.UserCreatePayload) (*model.User, error) {
	if _, err := mail.ParseAddress(p.Email); err != nil {
		return nil, ErrInvalidEmail
	}

	if len(p.Username) < 3 {
		return nil, ErrUsernameToShort
	}

	if utf8.RuneCountInString(p.Username) < 3 {
		return nil, ErrPasswordToShort
	}

	hash, err := password.Hash(p.Password)
	if err != nil {
		return nil, err
	}
	p.Password = *hash

	exist, err := u.user.ExistsByEmail(ctx, p.Email)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, ErrEmailAlreadyExists
	}

	exist, err = u.user.ExistsByUsername(ctx, p.Username)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, ErrUsernameAlreadyExists
	}

	id, err := u.user.Create(ctx, p)
	if err != nil {
		return nil, err
	}

	user, err := u.user.FindByID(ctx, *id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
