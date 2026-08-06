package repo

import (
	"context"
	"database/sql"
	"errors"
	"godo/internal/model"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
)

var (
	ErrSessionNotExist = errors.New("session not exists")
)

type Session struct {
	*Repo
	columns []string
	table   string
}

func (r *Repo) NewSession() *Session {
	return &Session{
		Repo:    r,
		table:   "sessions",
		columns: []string{"id", "csrf_token", "user_id", "expires_at"},
	}
}

func (s *Session) Create(ctx context.Context, p model.Session) error {
	_, err := s.sq.
		Insert(s.table).
		Columns(s.columns...).
		Values(p.ID, p.CSRFToken, p.UserID, p.ExpiresAt).
		ExecContext(ctx)
	return s.tErr(err)
}

func (s *Session) Find(ctx context.Context, sessionID string) (model.Session, error) {
	var r model.Session
	err := s.sq.
		Select(s.columns...).
		From(s.table).
		QueryRowContext(ctx).
		Scan(&r.ID, &r.CSRFToken, &r.UserID, &r.ExpiresAt)
	return r, s.tErr(err)
}

func (s *Session) Remove(ctx context.Context, sessionID string) error {
	_, err := s.sq.Delete(s.table).Where(sq.Eq{"id": sessionID}).ExecContext(ctx)
	return s.tErr(err)
}

func (s *Session) CleanUp(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			_, err := s.sq.Delete(s.table).Where(sq.Lt{"expires_at": time.Now()}).Exec()
			if err != nil {
				log.Println("cleanup session error:", err)
			}
		}
	}()
}

// =====: Helpers
// 1. translate error
func (s *Session) tErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrSessionNotExist
	}

	return err
}
