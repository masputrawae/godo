package model

import "time"

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

type UserCreatePayload struct {
	Username string
	Email    string
	Password string
}

type UserUpdatePayload struct {
	Username *string
	Email    *string
	Password *string
}

type UserLoginPayload struct {
	Username string
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
	Done       bool
	Due        *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserID     int
	StatusID   int
	PriorityID int
}

type TodoCreatePayload struct {
	UserID     int
	Task       string
	Due        *time.Time
	StatusID   *int
	PriorityID *int
}

type TodoUpdatePayload struct {
	UserID     int
	TodoID     int
	Task       *string
	Done       *bool
	Due        *time.Time
	StatusID   *int
	PriorityID *int
}

type Session struct {
	UserID    int
	CSRFToken string
	Expires   *time.Time
}
