package repo

import (
	"context"
	"database/sql"
	"godo/internal/model"

	sq "github.com/Masterminds/squirrel"
)

type Todo struct {
	*repo
}

func NewTodo(db *sql.DB) *Todo {
	return &Todo{New(db)}
}

func (t *Todo) Create(ctx context.Context, p model.TodoCreatePayload) (int, error) {
	q := t.sq.Insert("todos").Columns("task").Values(p.Task)

	if p.Due != nil {
		q = q.Columns("due").Values(*p.Due)
	}

	res, err := q.ExecContext(ctx)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func (t *Todo) FindAll(ctx context.Context) ([]model.Todo, error) {
	q := "SELECT id, task, done, due, created_at, updated_at, user_id, status_id, priority_id FROM todos"
	rows, err := t.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.Todo
	if rows.Next() {
		r, err := t.scan(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, rows.Err()
}

func (t *Todo) Update(ctx context.Context, todoID int, p model.TodoUpdatePayload) error {
	if p.Task == nil && p.Done == nil && p.Due == nil && p.StatusID == nil && p.PriorityID == nil {
		return nil
	}

	q := t.sq.Update("todos").Where(sq.Eq{"id": todoID}).Set("updated_at", "CURRENT_TIMESTAMP")

	if p.Task != nil {
		q = q.Set("task", *p.Task)
	}

	if p.Done != nil {
		q = q.Set("done", *p.Done)
	}

	if p.Due != nil {
		q = q.Set("due", *p.Due)
	}

	if p.PriorityID != nil {
		q = q.Set("priority_id", *p.PriorityID)
	}

	if p.StatusID != nil {
		q = q.Set("status_id", *p.StatusID)
	}

	_, err := q.ExecContext(ctx)
	return err
}

func (t *Todo) Delete(ctx context.Context, todoID int) error {
	_, err := t.db.ExecContext(ctx, "DELETE FROM todos WHERE id = ?", todoID)
	return err
}

func (t *Todo) scan(s interface{ Scan(dest ...any) error }) (model.Todo, error) {
	var r model.Todo
	err := s.Scan(&r.ID, &r.Task, &r.Done, &r.Due, &r.CreatedAt, &r.UpdatedAt, &r.UserID, &r.StatusID, &r.PriorityID)
	return r, err
}
