package repo

import (
	"context"
	"database/sql"
	"godo/internal/model"

	sq "github.com/Masterminds/squirrel"
)

type User struct {
	*repo
}

func NewUser(db *sql.DB) *User {
	return &User{New(db)}
}

func (u *User) Create(ctx context.Context, p model.UserCreatePayload) (int, error) {
	q := "INSERT INTO users (username, email, password) VALUES ( ? , ? , ? )"
	res, err := u.db.ExecContext(ctx, q, p.Username, p.Email, p.Password)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), err
}

func (u *User) FindByID(ctx context.Context, userID int) (model.User, error) {
	q := "SELECT username, email, password FROM users WHERE id = ?"
	return u.scan(u.db.QueryRowContext(ctx, q, userID))
}

func (u *User) FindByUsername(ctx context.Context, username string) (model.User, error) {
	q := "SELECT username, email, password FROM users WHERE username = ?"
	return u.scan(u.db.QueryRowContext(ctx, q, username))
}

func (u *User) Update(ctx context.Context, userID int, p model.UserUpdatePayload) error {
	if p.Password == nil && p.Username == nil && p.Email == nil {
		return nil
	}

	q := u.sq.
		Update("users").
		Where(sq.Eq{"id": userID}).
		Set("updated_at", "CURRENT_TIMESTAMP")

	if p.Password != nil {
		q = q.Set("password", *p.Password)
	}
	if p.Username != nil {
		q = q.Set("username", *p.Username)
	}
	if p.Email != nil {
		q = q.Set("email", *p.Email)
	}

	_, err := q.ExecContext(ctx)
	return err
}

func (u *User) Delete(ctx context.Context, userID int) error {
	_, err := u.db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", userID)
	return err
}

func (u *User) scan(s interface{ Scan(dest ...any) error }) (model.User, error) {
	var r model.User
	err := s.Scan(&r.Username, &r.Email, &r.Password)
	return r, err
}
