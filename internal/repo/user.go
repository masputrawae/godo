package repo

import (
	"context"
	"godo/internal/model"

	sq "github.com/Masterminds/squirrel"
)

type User struct {
	*repo
}

// =====: User
func NewUser(r *repo) *User {
	rp := &repo{
		db:    r.db,
		sq:    r.sq,
		table: "users",
		columns: []string{
			"id", "email", "username",
			"password", "created_at", "updated_at",
		},
	}
	return &User{rp}
}

// =====: Create
func (u *User) Create(ctx context.Context, p model.UserCreatePayload) (int, error) {
	res, err := u.sq.
		Insert(u.table).
		Columns("email", "username", "password").
		Values(p.Email, p.Username, p.Password).
		ExecContext(ctx)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
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

// =====: Find by Email
func (u *User) FindByEmail(ctx context.Context, email string) (model.User, error) {
	return u.selectOne(ctx, sq.Eq{"email": email})
}

// =====: Helpers
// 1. scan
func (u *User) scan(sc interface{ Scan(dest ...any) error }) (model.User, error) {
	var r model.User
	err := sc.Scan(&r.ID, &r.Email, &r.Username, &r.Password, &r.CreatedAt, &r.UpdatedAt)
	return r, err
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
