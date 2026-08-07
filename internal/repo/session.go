package repo

import (
	"context"
	"godo/internal/model"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Session struct {
	*repo
}

// =====: Session
func NewSession(r *repo) *Session {
	rp := &repo{
		db:      r.db,
		sq:      r.sq,
		table:   "sessions",
		columns: []string{"id", "user_id", "csrf_token", "expires_at"},
	}
	return &Session{rp}
}

// =====: Create
func (s *Session) Create(ctx context.Context, p model.Session) error {
	_, err := s.sq.
		Insert(s.table).
		Columns(s.columns...).
		Values(p.ID, p.UserID, p.CSRFToken, p.ExpiresAt).
		ExecContext(ctx)
	return err
}

// =====: Delete by ID
func (s *Session) DeleteByID(ctx context.Context, id string) error {
	return s.del(ctx, sq.Eq{"id": id})
}

// =====: Delete by Expires
func (s *Session) DeleteByExpires(ctx context.Context) error {
	return s.del(ctx, sq.Lt{"expires_at": time.Now()})
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

// =====: Helpers
// 1. delete
func (s *Session) del(ctx context.Context, args any) error {
	_, err := s.sq.
		Delete(s.table).
		Where(args).
		ExecContext(ctx)
	return err
}
