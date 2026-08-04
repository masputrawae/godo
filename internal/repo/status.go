package repo

import (
	"context"
	"database/sql"
	"godo/internal/model"
)

type Status struct {
	*repo
}

func NewStatus(db *sql.DB) *Status {
	return &Status{New(db)}
}

func (s *Status) FindAll(ctx context.Context) ([]model.Status, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, emoji, name FROM statuses")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []model.Status
	for rows.Next() {
		r, err := s.scan(rows)
		if err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	return results, rows.Err()
}

func (s *Status) FindByID(ctx context.Context, statusID int) (model.Status, error) {
	return s.scan(s.db.QueryRowContext(ctx, "SELECT id, emoji, name FROM statuses WHERE id = ?", statusID))
}

func (s *Status) scan(sc interface{ Scan(dest ...any) error }) (model.Status, error) {
	var r model.Status
	err := sc.Scan(&r.ID, &r.Emoji, &r.Name)
	return r, err
}
