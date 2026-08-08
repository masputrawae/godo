package repo

import (
	"context"
	"database/sql"
	"godo/internal/model"

	sq "github.com/Masterminds/squirrel"
)

type User struct {
	db      *sql.DB
	sq      sq.StatementBuilderType
	table   string
	columns []string
}

func NewUser(db *sql.DB) *User {
	return &User{
		db:      db,
		sq:      sq.StatementBuilder.PlaceholderFormat(sq.Question).RunWith(db),
		table:   "users",
		columns: []string{"id", "email", "username", "password", "created_at", "updated_at"},
	}
}

func (u *User) Create(ctx context.Context, p model.UserPayloadCreate) error {
	_, err := u.sq.
		Insert(u.table).
		Columns("email", "username", "password").
		Values(p.Email, p.Username, p.Password).
		ExecContext(ctx)
	return err
}

// =====: Find by ID
func (u *User) FindByID(ctx context.Context, id int) (model.User, error) {
	row := u.sq.
		Select(u.columns...).
		From(u.table).
		Where(sq.Eq{"id": id}).
		QueryRowContext(ctx)

	return u.scan(row)
}

// =====: Find by Username
func (u *User) FindByUsername(ctx context.Context, username string) (model.User, error) {
	row := u.sq.
		Select(u.columns...).
		From(u.table).
		Where(sq.Eq{"username": username}).
		QueryRowContext(ctx)

	return u.scan(row)
}

// =====: Delete by ID
func (u *User) DeleteByID(ctx context.Context, id int) error {
	_, err := u.sq.Delete(u.table).Where(sq.Eq{"id": id}).ExecContext(ctx)
	return err
}

// =====: Update
func (u *User) Update(ctx context.Context, p model.UserPayloadUpdate) error {
	q := u.sq.Update(u.table).Where(sq.Eq{"id": p.ID})

	if p.Email != nil {
		q = q.Set("email", *p.Email)
	}

	if p.Username != nil {
		q = q.Set("username", *p.Username)
	}

	if p.Password != nil {
		q = q.Set("password", *p.Password)
	}

	_, err := q.ExecContext(ctx)
	return err
}

// =====: Helpers
// 1. scan
func (u *User) scan(sc interface{ Scan(dest ...any) error }) (model.User, error) {
	var r model.User
	err := sc.Scan(&r.ID, &r.Email, &r.Username, &r.Password, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}
