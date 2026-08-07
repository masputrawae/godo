package model

import "time"

// =====: Todo
type Todo struct {
	ID         int
	Task       string
	Done       bool
	Due        *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserID     int
	StatusID   int
	PriorityID int
}

type TodoPayloadCreate struct {
	Task       string
	UserID     int
	Due        *time.Time
	StatusID   *int
	PriorityID *int
}

type TodoPayloadUpdate struct {
	ID         int
	UserID     int
	Done       *bool
	Task       *string
	Due        *time.Time
	StatusID   *int
	PriorityID *int
}

// =====: Status & Priority
type Status struct {
	ID    int
	Emoji string
	Name  string
}

type Priority struct {
	ID    int
	Emoji string
	Name  string
}
