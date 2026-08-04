package repo

import (
	"context"
	"database/sql"
	"godo/internal/model"
)

type Priority struct {
	*repo
}

func NewPriority(db *sql.DB) *Priority {
	return &Priority{New(db)}
}

func (s *Priority) FindAll(ctx context.Context) ([]model.Priority, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, emoji, name FROM priorities")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []model.Priority
	for rows.Next() {
		r, err := s.scan(rows)
		if err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	return results, rows.Err()
}

func (s *Priority) FindByID(ctx context.Context, priorityID int) (model.Priority, error) {
	return s.scan(s.db.QueryRowContext(ctx, "SELECT id, emoji, name FROM priorities WHERE id = ?", priorityID))
}

func (s *Priority) scan(sc interface{ Scan(dest ...any) error }) (model.Priority, error) {
	var r model.Priority
	err := sc.Scan(&r.ID, &r.Emoji, &r.Name)
	return r, err
}
