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
	ErrUsernameAlreadyExist = errors.New("username already exists")
	ErrEmailAlreadyExist    = errors.New("username already exists")
	ErrUserNotFound         = errors.New("user not found")
)

type User struct {
	*Repo
	columns []string
	table   string
}

// =====: New User
func (r *Repo) NewUser() *User {
	return &User{
		Repo:  r,
		table: "users",
		columns: []string{
			"id", "username", "password",
			"email", "created_at", "updated_at",
		},
	}
}

// =====: Create
func (u *User) Create(ctx context.Context, p model.UserCreatePayload) (int, error) {
	res, err := u.sq.
		Insert(u.table).
		Columns("username", "password", "email").
		Values(p.Username, p.Password, p.Email).
		ExecContext(ctx)

	if err != nil {
		return -1, u.tErr(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, u.tErr(err)
	}

	return int(id), nil
}

// =====: Find by ID
func (u *User) FindByID(ctx context.Context, id int) (model.User, error) {
	return u.selectOne(ctx, sq.Eq{"id": id})
}

// =====: Find by Username
func (u *User) FindByUsername(ctx context.Context, username string) (model.User, error) {
	return u.selectOne(ctx, sq.Eq{"username": username})
}

// =====: Helpers
// 1. scan
func (u *User) scan(sc interface{ Scan(dest ...any) error }) (model.User, error) {
	var r model.User
	err := sc.Scan(&r.ID, &r.Username, &r.Password, &r.Email, &r.CreatedAt, &r.UpdatedAt)
	return r, u.tErr(err)
}

// 2. select one
func (u *User) selectOne(ctx context.Context, eq sq.Eq) (model.User, error) {
	return u.scan(
		u.sq.
			Select(u.columns...).
			From(u.table).
			Where(eq).
			QueryRowContext(ctx),
	)
}

// 3. translate error
func (u *User) tErr(err error) error {
	if err == nil {
		return nil
	}
	se := strings.ToLower(err.Error())
	eq := func(q string) bool {
		return strings.Contains(se, q)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrUserNotFound
	}

	if eq("unique") && eq("users.username") {
		return ErrUsernameAlreadyExist
	}

	if eq("unique") && eq("users.email") {
		return ErrUsernameAlreadyExist
	}

	return err
}
