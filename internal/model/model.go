package model

import "time"

type User struct {
	ID       int
	Username string
	Email    string
	Password string
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

type Todo struct {
	ID         int
	Task       string
	Done       *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserID     int
	StatusID   int
	PriorityID int
}
