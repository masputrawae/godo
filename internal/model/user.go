package model

import "time"

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
