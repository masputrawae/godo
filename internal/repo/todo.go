package repo

import (
	"context"
	"database/sql"
	"godo/internal/model"

	sq "github.com/Masterminds/squirrel"
)

type Todo struct {
	db      *sql.DB
	sq      sq.StatementBuilderType
	table   string
	columns []string
}

func NewTodo(db *sql.DB) *Todo {
	return &Todo{
		db:    db,
		sq:    sq.StatementBuilder.PlaceholderFormat(sq.Question).RunWith(db),
		table: "todos",
		columns: []string{
			"id", "task", "done",
			"due", "created_at", "updated_at",
			"user_id", "status_id", "priority_id",
		},
	}
}

// =====: Create
func (t *Todo) Create(ctx context.Context, p model.TodoPayloadCreate) error {
	q := t.sq.
		Insert(t.table).
		Columns("task", "user_id").
		Values(p.Task, p.UserID)

	if p.Due != nil {
		q = q.Columns("due").Values(*p.Due)
	}

	if p.StatusID != nil {
		q = q.Columns("status_id").Values(*p.StatusID)
	}

	if p.PriorityID != nil {
		q = q.Columns("priority_id").Values(*p.PriorityID)
	}

	_, err := q.ExecContext(ctx)
	return err
}

// =====: Update
func (t *Todo) Update(ctx context.Context, p model.TodoPayloadUpdate) error {
	if p.Done == nil && p.Task == nil && p.Due == nil && p.StatusID == nil && p.PriorityID == nil {
		return nil
	}

	q := t.sq.
		Update(t.table).
		Set("updated_at", "CURRENT_TIMESTAMP").
		Where(sq.And{
			sq.Eq{"user_id": p.UserID},
			sq.Eq{"id": p.ID},
		})

	if p.Done != nil {
		q = q.Set("done", *p.Done)
	}

	if p.Task != nil {
		q = q.Set("task", *p.Task)
	}

	if p.Due != nil {
		q = q.Set("due", *p.Due)
	}

	if p.StatusID != nil {
		q = q.Set("status_id", *p.StatusID)
	}

	if p.PriorityID != nil {
		q = q.Set("priority_id", *p.PriorityID)
	}

	_, err := q.ExecContext(ctx)
	return err
}

// =====: Find All
func (t *Todo) FindAll(ctx context.Context, userID int) ([]model.Todo, error) {
	rows, err := t.sq.
		Select(t.columns...).
		Where(sq.Eq{"user_id": userID}).
		OrderBy("priority_id ASC", "due DESC", "done ASC").
		QueryContext(ctx)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []model.Todo
	for rows.Next() {
		r, err := t.scan(rows)
		if err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	if rows.Err() == nil {
		return nil, err
	}

	return results, nil
}

// =====: Find by ID
func (t *Todo) FindByID(ctx context.Context, userID, id int) (model.Todo, error) {
	return t.scan(
		t.sq.
			Select(t.columns...).
			From(t.table).
			Where(sq.And{
				sq.Eq{"user_id": userID},
				sq.Eq{"id": id},
			}).
			QueryRowContext(ctx),
	)
}

// =====: Delete by ID
func (t *Todo) DeleteByID(ctx context.Context, userID, id int) error {
	_, err := t.sq.
		Delete(t.table).
		Where(sq.And{
			sq.Eq{"user_id": userID},
			sq.Eq{"id": id},
		}).
		ExecContext(ctx)

	return err
}

// =====: Helpers
// 1. scan
func (t *Todo) scan(sc interface{ Scan(dest ...any) error }) (model.Todo, error) {
	var r model.Todo
	err := sc.Scan(
		&r.ID, &r.Task, &r.Done,
		&r.Due, &r.CreatedAt, &r.UpdatedAt,
		&r.UserID, &r.StatusID, &r.PriorityID,
	)
	return r, err
}
