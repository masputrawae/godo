package repo

import (
	"context"
	"database/sql"
	"godo/internal/model"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Session struct {
	db      *sql.DB
	sq      sq.StatementBuilderType
	table   string
	columns []string
}

func NewSession(db *sql.DB) *Session {
	return &Session{
		db:      db,
		sq:      sq.StatementBuilder.PlaceholderFormat(sq.Question).RunWith(db),
		table:   "sessions",
		columns: []string{"id", "user_id", "csrf_token", "expires_at"},
	}
}

// =====: Create
func (s *Session) Create(ctx context.Context, p model.Session) error {
	_, err := s.sq.
		Insert(s.table).
		Columns("id", "user_id", "csrf_token", "expires_at").
		Values(p.ID, p.UserID, p.CSRFToken, p.ExpiresAt).
		ExecContext(ctx)
	return err
}

// =====: Delete by ID
func (s *Session) DeleteByID(ctx context.Context, id string) error {
	_, err := s.sq.
		Delete(s.table).
		Where(sq.Eq{"id": id}).
		ExecContext(ctx)
	return err
}

// =====: Delete by Expires
func (s *Session) DeleteByExpires(ctx context.Context) error {
	_, err := s.sq.
		Delete(s.table).
		Where(sq.LtOrEq{"expires_at": time.Now()}).
		ExecContext(ctx)
	return err
}

// =====: Find by ID
func (s *Session) FindByID(ctx context.Context, id string) (model.Session, error) {
	var r model.Session
	err := s.sq.
		Select(s.columns...).
		From(s.table).
		Where(sq.Eq{"id": id}).
		QueryRowContext(ctx).
		Scan(&r.ID, &r.UserID, &r.CSRFToken, &r.ExpiresAt)
	return r, err
}
