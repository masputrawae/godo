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

type UserPayloadCreate struct {
	Email    string
	Username string
	Password string
}

type UserPayloadLogin struct {
	Username string
	Password string
}

type UserPayloadUpdate struct {
	ID       int
	Email    *string
	Username *string
	Password *string
}
