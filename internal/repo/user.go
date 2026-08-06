package repo

import (
	"context"
	"database/sql"
	"errors"
	"godo/internal/model"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

var (
	ErrUniqueUsername  = errors.New("unique username")
	ErrUniqueEmail     = errors.New("unique email")
	ErrNotNullUsername = errors.New("not null username")
	ErrNotNullPassword = errors.New("not null password")
	ErrNotNullEmail    = errors.New("not null email")
	ErrUserNotFound    = errors.New("user not found")
)

type User struct {
	*Repo
}

func (u *Repo) NewUser() *User {
	repo := &Repo{
		db:    u.db,
		sq:    u.sq,
		table: "users",
		columns: []string{
			"id", "username", "password",
			"email", "created_at", "updated_at",
		},
	}
	return &User{repo}
}

// =====: CREATE USER
func (u *User) Create(ctx context.Context, p model.UserCreatePayload) (*int, error) {
	res, err := u.sq.
		Insert(u.table).
		Columns("username", "password", "email").
		Values(p.Username, p.Password, p.Email).
		ExecContext(ctx)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return new(int(id)), nil
}

// =====: FIND USER BY ID
func (u *User) FindByID(ctx context.Context, id int) (*model.User, error) {
	return u.selectOne(ctx, sq.Eq{"id": id})
}

// =====: FIND USER BY USERNAME
func (u *User) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	return u.selectOne(ctx, sq.Eq{"username": username})
}

// =====: FIND USER BY EMAIL
func (u *User) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return u.selectOne(ctx, sq.Eq{"email": email})
}

// =====: HELPERS
// 1. scan
func (u *User) scan(sc interface{ Scan(dest ...any) error }) (*model.User, error) {
	var res model.User
	err := sc.Scan(&res.ID, &res.Username, &res.Password, &res.Email, &res.CreatedAt, &res.UpdatedAt)
	return &res, err
}

// 2. select one
func (u *User) selectOne(ctx context.Context, eq sq.Eq) (*model.User, error) {
	user, err := u.scan(
		u.sq.
			Select(u.columns...).
			From(u.table).
			Where(eq).
			QueryRowContext(ctx),
	)
	return user, err
}

// =====: TRANSLATE ERRORS
func (r *User) TranslateError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrUserNotFound
	}

	se := strings.ToLower(err.Error())

	if strings.Contains(se, "unique") {
		switch {
		case strings.Contains(se, "users.username"):
			return ErrUniqueUsername
		case strings.Contains(se, "users.email"):
			return ErrUniqueEmail
		}
	}

	if strings.Contains(se, "not null") {
		switch {
		case strings.Contains(se, "users.username"):
			return ErrNotNullUsername
		case strings.Contains(se, "users.password"):
			return ErrNotNullPassword
		case strings.Contains(se, "users.email"):
			return ErrNotNullEmail
		}
	}

	return err
}
