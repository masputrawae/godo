package service

import (
	"context"
	"database/sql"
	"errors"
	"godo/internal/model"
	"godo/internal/repo"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameOrPasswordNotMatch = errors.New("username or password not match")
	ErrPasswordToShort            = errors.New("password to short. min 8 character")
	ErrUsernameToShort            = errors.New("username to short. min 3 character")
	ErrUsernameAlreadyUse         = errors.New("username already in use")
	ErrEmailAlreadyUse            = errors.New("email already in use")
	ErrEmailInvalid               = errors.New("email invalid")
	ErrUserAccountNotFound        = errors.New("user account not found")
	ErrPasswordCannotSame         = errors.New("passwords cannot be the same")
	ErrEmailCannotSame            = errors.New("email cannot be the same")
	ErrUsernameCannotSame         = errors.New("username cannot be the same")
	ErrNothingHasChanged          = errors.New("nothing has changed")
)

type User struct {
	repo *repo.User
}

func NewUser(repo *repo.User) *User {
	return &User{repo: repo}
}

func (u *User) Login(ctx context.Context, p model.UserLoginPayload) (*model.User, error) {
	user, err := u.GetByUsername(ctx, p.Username)
	if err != nil {
		return nil, err
	}

	if !checkPassword(user.Password, p.Password) {
		return nil, ErrUsernameOrPasswordNotMatch
	}

	// avoid password leaks
	user.Password = ""
	return user, nil
}

func (u *User) Register(ctx context.Context, p model.UserCreatePayload) (*model.User, error) {
	if len(p.Password) < 8 {
		return nil, ErrPasswordToShort
	}

	if len(p.Username) < 3 {
		return nil, ErrUsernameToShort
	}

	if _, err := mail.ParseAddress(p.Email); err != nil {
		return nil, ErrEmailInvalid
	}

	hash, err := hashPassword(p.Password)
	if err != nil {
		return nil, err
	}
	p.Password = hash

	id, err := u.repo.Create(ctx, p)
	if err != nil {
		if field, ok := uniqueViolationField(err); ok {
			switch field {
			case "email":
				return nil, ErrEmailAlreadyUse
			case "username":
				return nil, ErrUsernameAlreadyUse
			default:
				return nil, err
			}
		}

		return nil, err
	}

	user, err := u.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (u *User) Update(ctx context.Context, userID int, p model.UserUpdatePayload) (*model.User, error) {
	if p.Email == nil && p.Username == nil && p.Password == nil {
		return nil, ErrNothingHasChanged
	}

	old, err := u.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if p.Password != nil {
		if len(*p.Password) < 8 {
			return nil, ErrPasswordToShort
		}

		if checkPassword(old.Password, *p.Password) {
			return nil, ErrPasswordCannotSame
		}

		hash, err := hashPassword(*p.Password)
		if err != nil {
			return nil, err
		}

		p.Password = &hash
	}

	if p.Username != nil && len(*p.Username) < 3 {
		return nil, ErrUsernameToShort
	}

	if p.Email != nil {
		if _, err := mail.ParseAddress(*p.Email); err != nil {
			return nil, ErrEmailInvalid
		}
	}

	if err := u.repo.Update(ctx, userID, p); err != nil {
		if field, ok := uniqueViolationField(err); ok {
			switch field {
			case "email":
				return nil, ErrEmailAlreadyUse
			case "username":
				return nil, ErrUsernameAlreadyUse
			default:
				return nil, err
			}
		}

		return nil, err
	}

	return u.GetByID(ctx, userID)
}

func (u *User) GetByID(ctx context.Context, userID int) (*model.User, error) {
	user, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserAccountNotFound
		}
	}

	user.Password = ""
	return &user, nil
}

func (u *User) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := u.repo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserAccountNotFound
		}
	}

	user.Password = ""
	return &user, nil
}

// hash password
func hashPassword(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// check password
func checkPassword(h, p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p)) == nil
}

// check duplication
func uniqueViolationField(err error) (string, bool) {
	if err == nil {
		return "", false
	}
	msg := err.Error()
	const prefix = "UNIQUE constraint failed:"
	if !strings.Contains(msg, prefix) {
		return "", false
	}
	after := strings.TrimSpace(strings.TrimPrefix(msg, prefix))
	if strings.Contains(after, ".") {
		parts := strings.Split(after, ".")
		return parts[len(parts)-1], true
	}
	return after, true
}
