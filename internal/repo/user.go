package repo

import (
	"context"
	"database/sql"
	"errors"
	"godo/internal/model"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
)

var (
	ErrDataNotFound = errors.New("data not found")
)

type User struct {
	*Repo
}

func (u *Repo) NewUser() *User {
	repo := &Repo{
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
	slog.Info("⏳ create user")

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

	slog.Info("✅ successfully saved to the database")
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

// =====: CHECKING DUPLICATION DATA
func (u *User) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := u.db.
		QueryRowContext(
			ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email,
		).
		Scan(&exists)
	return exists, err
}

func (u *User) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := u.db.
		QueryRowContext(
			ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username,
		).
		Scan(&exists)
	return exists, err
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

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrDataNotFound
	}

	return user, err
}
