package service

import (
	"context"
	"godo/internal/model"
	"godo/internal/repo"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	session *repo.Session
	ttl     time.Duration
}

// =====: Session
func NewSession(session *repo.Session) *Session {
	return &Session{
		session: session,
		ttl:     1 * time.Minute,
	}
}

// =====: Set Session
func (s *Session) Set(ctx context.Context, w http.ResponseWriter, userID int) error {
	uID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	sessionID := uID.String()

	uID, err = uuid.NewV7()
	if err != nil {
		return err
	}

	csrfToken := uID.String()

	expires := time.Now().Add(s.ttl)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		Expires:  expires,
		MaxAge:   int(s.ttl.Seconds()),
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return s.session.Create(ctx, model.Session{
		ID:        sessionID,
		CSRFToken: csrfToken,
		ExpiresAt: expires,
		UserID:    userID,
	})
}

// =====: Get Session
func (s *Session) Get(ctx context.Context, sessionID string) (model.Session, error) {
	return s.session.FindByID(ctx, sessionID)
}

// =====: Delete Session
func (s *Session) Delete(ctx context.Context, w http.ResponseWriter, sessionID string) error {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	return s.session.DeleteByID(ctx, sessionID)
}
