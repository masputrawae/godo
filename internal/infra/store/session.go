package store

import (
	"godo/internal/model"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

var sessions = make(map[string]model.Session)

type Session struct {
	sessions map[string]model.Session
	expires  time.Time
	maxAge   time.Duration
	mu       sync.RWMutex
}

func New() *Session {
	return &Session{
		sessions: sessions,
	}
}

func (s *Session) SaveSession(userID int, sessionID, csrfToken string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[sessionID] = model.Session{
		UserID:    userID,
		CSRFToken: csrfToken,
		ExpiresAt: s.expires,
	}
}

func (s *Session) DeleteSession(sessionID string) {
	s.mu.Lock()
	delete(s.sessions, sessionID)
	s.mu.Unlock()
}

func (s *Session) GetSession(sessionID string) (model.Session, bool) {
	s.mu.RLock()
	session, exist := s.sessions[sessionID]
	s.mu.RUnlock()

	return session, exist
}

func (s *Session) SetCookie(w http.ResponseWriter, sessionID string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		Expires:  s.expires,
		MaxAge:   int(s.maxAge),
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (s *Session) DeleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

func (s *Session) ParseCTX(r *http.Request) (model.Session, bool) {
	session, ok := r.Context().Value("session").(model.Session)
	return session, ok
}

func (s *Session) GenID() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	return id.String(), nil
}
