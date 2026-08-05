package session

import (
	"godo/internal/model"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	sessions map[string]model.Session
	mu       sync.RWMutex
	maxAge   int
	expires  time.Time
}

func New() *Session {
	return &Session{
		sessions: make(map[string]model.Session),
		maxAge:   int(24 * time.Hour),
		expires:  time.Now().Add(24 * time.Hour),
	}
}

func (s *Session) Set(w http.ResponseWriter, userID int) error {
	sessionID, err := genID()
	if err != nil {
		return err
	}
	csrfToken, err := genID()
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sessionID] = model.Session{
		UserID:    userID,
		CSRFToken: csrfToken,
		Expires:   &s.expires,
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Expires:  s.expires,
		Path:     "/",
		MaxAge:   s.maxAge,
		Secure:   false, // dev only
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

func (s *Session) Get(w http.ResponseWriter, r *http.Request) *model.Session {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil
	}

	id := cookie.Value

	s.mu.RLock()
	session, ok := s.sessions[id]
	s.mu.RUnlock()

	if !ok {
		return nil
	}

	if session.Expires.Before(time.Now()) {
		s.Delete(w, id)
		return nil
	}

	return &session
}

func (s *Session) Delete(w http.ResponseWriter, sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionID)

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (s *Session) Parse(r *http.Request) *model.Session {
	session, ok := r.Context().Value("session").(*model.Session)
	if !ok {
		return nil
	}

	return session
}

func genID() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
