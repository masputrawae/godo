package model

import "time"

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

type TodoCreatePayload struct {
	Task       string
	Due        *time.Time
	StatusID   *int
	PriorityID *int
}

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
