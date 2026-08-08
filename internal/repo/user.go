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
	ErrUserNotFound        = errors.New("user not found")
	ErrEmailAlreadyUsed    = errors.New("email already used")
	ErrUsernameAlreadyUsed = errors.New("username already used")
	ErrUserAlreadyUsed     = errors.New("user already used")
)

type User struct {
	sq      sq.StatementBuilderType
	table   string
	columns []string
}

func NewUser(db *sql.DB) *User {
	return &User{
		sq:    sq.StatementBuilder.PlaceholderFormat(sq.Question).RunWith(db),
		table: "users",
		columns: []string{
			"id", "email", "username",
			"password", "created_at", "updated_at",
		},
	}
}

// create new user
func (u *User) Create(ctx context.Context, p model.UserRequestRegister) (*model.User, error) {
	row := u.sq.
		Insert(u.table).
		Columns("email", "username", "password").
		Values(p.Email, p.Username, p.Password).
		Suffix("RETURNING " + strings.Join(u.columns, ",")).
		QueryRowContext(ctx)

	return u.scan(row)
}

// update user
func (u *User) Update(ctx context.Context, id int, p model.UserRequestUpdate) (*model.User, error) {
	if p.Password == nil && p.Username == nil && p.Email == nil {
		return nil, nil
	}

	q := u.sq.
		Update(u.table).
		Where(sq.Eq{"id": id}).
		Set("updated_at", sq.Expr("CURRENT_TIMESTAMP"))

	if p.Password != nil {
		q = q.Set("password", *p.Password)
	}
	if p.Username != nil {
		q = q.Set("username", *p.Username)
	}
	if p.Email != nil {
		q = q.Set("email", *p.Email)
	}

	row := q.
		Suffix("RETURNING " + strings.Join(u.columns, ",")).
		QueryRowContext(ctx)

	return u.scan(row)
}

// delete user by user id
func (u *User) DeleteByID(ctx context.Context, id int) error {
	_, err := u.sq.Delete(u.table).Where(sq.Eq{"id": id}).ExecContext(ctx)
	return err
}

// find user by id
func (u *User) FindByID(ctx context.Context, id int) (*model.User, error) {
	return u.findOne(ctx, sq.Eq{"id": id})
}

// find user by username
func (u *User) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	return u.findOne(ctx, sq.Eq{"username": username})
}

// find one
func (u *User) findOne(ctx context.Context, eq sq.Eq) (*model.User, error) {
	row := u.sq.
		Select(u.columns...).
		From(u.table).
		Where(eq).
		QueryRowContext(ctx)

	return u.scan(row)
}

// scan
func (u *User) scan(i interface{ Scan(dest ...any) error }) (*model.User, error) {
	var r = new(model.User)
	err := i.Scan(
		&r.ID, &r.Email, &r.Username,
		&r.Password, &r.CreatedAt, &r.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		msg := err.Error()

		switch {
		case strings.Contains(msg, "UNIQUE constraint failed: users.email"):
			return nil, ErrEmailAlreadyUsed
		case strings.Contains(msg, "UNIQUE constraint failed: users.username"):
			return nil, ErrUsernameAlreadyUsed
		}

		return nil, err
	}

	return r, err
}
