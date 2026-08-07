package service

import (
	"context"
	"godo/internal/model"
	"godo/internal/repo"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	rpSession  *repo.Session
	TTL        time.Duration
	CookieName string
}

func NewSession(rpSession *repo.Session, TTL time.Duration, CookieName string) *Session {
	return &Session{rpSession, TTL, CookieName}
}

// =====: Set
func (s *Session) Set(ctx context.Context, w http.ResponseWriter, userID int) error {
	expires := time.Now().Add(s.TTL)
	maxAge := int(s.TTL.Seconds())

	sessionID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	csrfToken, err := uuid.NewV7()
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     s.CookieName,
		Value:    sessionID.String(),
		Path:     "/",
		Expires:  expires,
		MaxAge:   maxAge,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return s.rpSession.Create(ctx, model.Session{
		ID:        sessionID.String(),
		CSRFToken: csrfToken.String(),
		UserID:    userID,
		ExpiresAt: expires,
	})
}

// =====: Get
func (s *Session) Get(ctx context.Context, sessionID string) (*model.Session, error) {
	session, err := s.rpSession.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, s.rpSession.DeleteByExpires(ctx)
	}

	return &session, nil
}

// =====: Parse Ctx
func (s *Session) ParseCtx(r *http.Request) (*model.Session, bool) {
	session, ok := r.Context().Value(s.CookieName).(model.Session)
	return &session, ok
}

// =====: Remove
func (s *Session) Remove(ctx context.Context, w http.ResponseWriter, sessionID string) error {
	http.SetCookie(w, &http.Cookie{
		Name:   s.CookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	return s.rpSession.DeleteByID(ctx, sessionID)
}

// =====: Auto Clean
func (s *Session) AutoClean(ctx context.Context, ticker *time.Ticker) {
	if ticker == nil {
		ticker = time.NewTicker(s.TTL)
	}

	for {
		select {
		case <-ticker.C:
			if err := s.rpSession.DeleteByExpires(ctx); err != nil {
				log.Println("auto clean session error:", err)
			}
		case <-ctx.Done():
			log.Println("auto clean session stop")
		}
	}
}
