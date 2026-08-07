package model

import "time"

type Session struct {
	ID        string
	UserID    int
	CSRFToken string
	ExpiresAt time.Time
}
