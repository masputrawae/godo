package store

import (
	"godo/internal/model"
	"sync"
	"time"
)

var sessions = make(map[string]model.Session)

type Session struct {
	sessions map[string]model.Session
	expires  time.Time
	mu       sync.RWMutex
}

func NewSession() *Session {
	return &Session{
		sessions: sessions,
		expires:  time.Now().Add(24 * time.Hour),
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
