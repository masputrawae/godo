package repo

import (
	"context"
	"database/sql"
	"godo/internal/model"

	sq "github.com/Masterminds/squirrel"
)

type Priority struct {
	db      *sql.DB
	sq      sq.StatementBuilderType
	table   string
	columns []string
}

func NewPriority(db *sql.DB) *Priority {
	return &Priority{
		db:      db,
		sq:      sq.StatementBuilder.PlaceholderFormat(sq.Question).RunWith(db),
		table:   "priorities",
		columns: []string{"id", "emoji", "name"},
	}
}

// =====: Find All
func (s *Priority) FindAll(ctx context.Context) ([]model.Priority, error) {
	rows, err := s.sq.
		Select(s.columns...).
		From(s.table).
		OrderBy("id ASC").
		QueryContext(ctx)

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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// =====: Find by ID
func (s *Priority) FindByID(ctx context.Context, id int) (model.Priority, error) {
	row := s.sq.
		Select(s.columns...).
		From(s.table).
		Where(sq.Eq{"id": id}).
		QueryRowContext(ctx)

	return s.scan(row)
}

// =====: Helpers
// 1. scan
func (s *Priority) scan(sc interface{ Scan(dest ...any) error }) (model.Priority, error) {
	var r model.Priority
	err := sc.Scan(&r.ID, &r.Emoji, &r.Name)
	return r, err
}
