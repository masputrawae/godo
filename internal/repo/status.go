package repo

import (
	"context"
	"database/sql"
	"godo/internal/model"

	sq "github.com/Masterminds/squirrel"
)

type Status struct {
	db      *sql.DB
	sq      sq.StatementBuilderType
	table   string
	columns []string
}

func NewStatus(db *sql.DB) *Status {
	return &Status{
		db:      db,
		sq:      sq.StatementBuilder.PlaceholderFormat(sq.Question).RunWith(db),
		table:   "statuses",
		columns: []string{"id", "emoji", "name"},
	}
}

// =====: Find All
func (s *Status) FindAll(ctx context.Context) ([]model.Status, error) {
	rows, err := s.sq.
		Select(s.columns...).
		From(s.table).
		OrderBy("id ASC").
		QueryContext(ctx)

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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// =====: Find by ID
func (s *Status) FindByID(ctx context.Context, id int) (model.Status, error) {
	row := s.sq.
		Select(s.columns...).
		From(s.table).
		Where(sq.Eq{"id": id}).
		QueryRowContext(ctx)

	return s.scan(row)
}

// =====: Helpers
// 1. scan
func (s *Status) scan(sc interface{ Scan(dest ...any) error }) (model.Status, error) {
	var r model.Status
	err := sc.Scan(&r.ID, &r.Emoji, &r.Name)
	return r, err
}
