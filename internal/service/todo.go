package service

import (
	"context"
	"database/sql"
	"errors"
	"godo/internal/model"
	"godo/internal/repo"
)

var (
	ErrInvalidIDPriority = errors.New("Invalid ID Priority")
	ErrInvalidIDStatus   = errors.New("Invalid ID Status")
	ErrTodoNotFound      = errors.New("Todo not found")
)

type Todo struct {
	rpTodo     *repo.Todo
	rpPriority *repo.Priority
	rpStatus   *repo.Status
}

func NewTodo(rpTodo *repo.Todo, rpPriority *repo.Priority, rpStatus *repo.Status) *Todo {
	return &Todo{rpTodo, rpPriority, rpStatus}
}

func (t *Todo) Create(ctx context.Context, p model.TodoPayloadCreate) error {
	if p.PriorityID != nil {
		_, err := t.rpPriority.FindByID(ctx, *p.PriorityID)
		if err != nil {
			return ErrInvalidIDPriority
		}
	}

	if p.StatusID != nil {
		_, err := t.rpStatus.FindByID(ctx, *p.PriorityID)
		if err != nil {
			return ErrInvalidIDPriority
		}
	}

	return t.rpTodo.Create(ctx, p)
}

func (t *Todo) Update(ctx context.Context, p model.TodoPayloadUpdate) error {
	old, err := t.rpTodo.FindByID(ctx, p.UserID, p.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTodoNotFound
		}
		return err
	}

	if old.Task == *p.Task {
		p.Task = nil
	}

	if old.Due == p.Due {
		p.Due = nil
	}

	if old.Done == *p.Done {
		p.Done = nil
	}

	if p.PriorityID != nil && p.PriorityID != &old.PriorityID {
		_, err := t.rpPriority.FindByID(ctx, *p.PriorityID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrInvalidIDPriority
			}
			return err
		}
	}

	if p.StatusID != nil && p.StatusID != &old.StatusID {
		_, err := t.rpStatus.FindByID(ctx, *p.PriorityID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrInvalidIDPriority
			}
			return err
		}
	}

	return t.rpTodo.Update(ctx, p)
}

func (t *Todo) FindAll(ctx context.Context, userID int) ([]model.Todo, error) {
	return t.rpTodo.FindAll(ctx, userID)
}

func (t *Todo) Delete(ctx context.Context, userID, id int) error {
	return t.rpTodo.DeleteByID(ctx, userID, id)
}
