package model

import "time"

type User struct {
	ID        int
	Email     string
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRequestLogin struct {
	Username string `validate:"required,min=1,max=32"`
	Password string `validate:"required,min=1,max=255"`
}

type UserRequestRegister struct {
	Email    string `validate:"required,email"`
	Username string `validate:"required,min=3,max=32"`
	Password string `validate:"required,min=8,max=255"`
}

type UserRequestUpdate struct {
	Email    *string `validate:"omitempty,email"`
	Username *string `validate:"omitempty,min=3,max=32"`
	Password *string `validate:"omitempty,min=8,max=255"`
}
