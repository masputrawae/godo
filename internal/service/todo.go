package service

import (
	"context"
	"database/sql"
	"errors"
	"godo/internal/model"
	"godo/internal/repo"
)

var (
	ErrTodoNotFound               = errors.New("todo not found")
	ErrTodoListStillEmpty         = errors.New("the task list is still empty")
	ErrStatusOptionNotAvailable   = errors.New("status option not available")
	ErrPriorityOptionNotAvailable = errors.New("status option not available")
)

type Todo struct {
	todoRepo     *repo.Todo
	statusRepo   *repo.Status
	priorityRepo *repo.Priority
}

func NewTodo(t *repo.Todo, s *repo.Status, p *repo.Priority) *Todo {
	return &Todo{todoRepo: t, statusRepo: s, priorityRepo: p}
}

func (t *Todo) Add(ctx context.Context, p model.TodoCreatePayload) (*model.Todo, error) {
	if p.StatusID != nil {
		_, err := t.GetStatusByID(ctx, *p.StatusID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrStatusOptionNotAvailable
			}
			return nil, err
		}
	}

	if p.PriorityID != nil {
		_, err := t.GetStatusByID(ctx, *p.StatusID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrPriorityOptionNotAvailable
			}
			return nil, err
		}
	}

	id, err := t.todoRepo.Create(ctx, p)
	if err != nil {
		return nil, err
	}

	return t.GetByID(ctx, p.UserID, id)
}

func (t *Todo) GetAll(ctx context.Context, userID int) ([]model.Todo, error) {
	todos, err := t.todoRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(todos) == 0 {
		return nil, ErrTodoListStillEmpty
	}
	return todos, nil
}

func (t *Todo) Delete(ctx context.Context, userID, todoID int) error {
	return t.todoRepo.Delete(ctx, userID, todoID)
}

func (t *Todo) Update(ctx context.Context, p model.TodoUpdatePayload) (*model.Todo, error) {
	if p.Task == nil && p.StatusID == nil && p.PriorityID == nil && p.Done == nil && p.Due == nil {
		return nil, ErrNothingHasChanged
	}

	old, err := t.GetByID(ctx, p.UserID, p.TodoID)
	if err != nil {
		return nil, err
	}

	if p.Task == nil || old.Task == *p.Task {
		p.Task = nil
	}

	if p.StatusID == nil || old.StatusID == *p.StatusID {
		p.StatusID = nil
	} else {
		if _, err := t.GetStatusByID(ctx, *p.StatusID); err != nil {
			return nil, err
		}
	}

	if p.PriorityID == nil || old.PriorityID == *p.PriorityID {
		p.PriorityID = nil
	} else {
		if _, err := t.GetPriorityByID(ctx, *p.PriorityID); err != nil {
			return nil, err
		}
	}

	if p.Due == nil && old.Due == p.Due {
		p.Task = nil
	}

	if err = t.todoRepo.Update(ctx, p); err != nil {
		return nil, err
	}

	return t.GetByID(ctx, p.UserID, p.TodoID)
}

func (t *Todo) GetByID(ctx context.Context, userID, todoID int) (*model.Todo, error) {
	todo, err := t.todoRepo.FindByID(ctx, userID, todoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTodoNotFound
		}
		return nil, err
	}
	return &todo, nil
}

func (t *Todo) GetAllStatuses(ctx context.Context) ([]model.Status, error) {
	return t.statusRepo.FindAll(ctx)
}

func (t *Todo) GetAllPriorities(ctx context.Context) ([]model.Priority, error) {
	return t.priorityRepo.FindAll(ctx)
}

func (t *Todo) GetStatusByID(ctx context.Context, statusID int) (*model.Status, error) {
	status, err := t.statusRepo.FindByID(ctx, statusID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPriorityOptionNotAvailable
		}
		return nil, err
	}
	return &status, nil
}

func (t *Todo) GetPriorityByID(ctx context.Context, priorityID int) (*model.Priority, error) {
	priority, err := t.priorityRepo.FindByID(ctx, priorityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPriorityOptionNotAvailable
		}
		return nil, err
	}
	return &priority, nil
}
