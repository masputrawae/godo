package model

import "time"

// =====: User
type User struct {
	ID        int
	Email     string
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserCreatePayload struct {
	Email    string
	Username string
	Password string
}

type UserLoginPayload struct {
	Username string
	Password string
}

// =====: Session
type Session struct {
	ID        string
	UserID    int
	CSRFToken string
	ExpiresAt time.Time
}
