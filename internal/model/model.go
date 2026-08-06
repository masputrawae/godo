package model

import "time"

// =====: Session Model
type Session struct {
	ID        string
	CSRFToken string
	UserID    int
	ExpiresAt time.Time
}

// =====: User Model
type User struct {
	ID        int
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserCreatePayload struct {
	Username string
	Password string
	Email    string
}

type UserLoginPayload struct {
	Username string
	Password string
}

// =====: Todo Model
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
	UserId     int
	Task       string
	Due        *time.Time
	StatusID   *int
	PriorityID *int
}

type TodoUpdatePayload struct {
	UserId     int
	Task       *string
	Due        *time.Time
	StatusID   *int
	PriorityID *int
}
