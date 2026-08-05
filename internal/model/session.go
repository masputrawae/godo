package model

import "time"

type Session struct {
	UserID    int
	CSRFToken string
	ExpiresAt time.Time
}
