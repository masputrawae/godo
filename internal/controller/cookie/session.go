package cookie

import (
	"godo/internal/infra/store"
	"godo/internal/model"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	store   *store.Session
	name    string
	expires time.Time
	maxAge  int
}

func NewSession(store *store.Session) *Session {
	return &Session{
		store:   store,
		name:    "session",
		expires: time.Now().Add(24 * time.Hour),
		maxAge:  int(24 * time.Hour),
	}
}

func (s *Session) Set(w http.ResponseWriter, userID int, sessionID, csrfToken string) error {
	s.store.SaveSession(userID, sessionID, csrfToken)

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		Expires:  s.expires,
		MaxAge:   s.maxAge,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

func (s *Session) Get(r *http.Request) (*model.Session, bool) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil, false
	}

	session, ok := s.store.GetSession(cookie.Value)
	if !ok {
		return nil, false
	}

	return &session, true
}

func (s *Session) Delete(w http.ResponseWriter, sessionID string) {
	s.store.DeleteSession(sessionID)
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

func (s *Session) Parse(r *http.Request) (*model.Session, bool) {
	session, ok := r.Context().Value("session").(model.Session)
	if !ok {
		return nil, false
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, false
	}

	return &session, true
}

func (s *Session) GenID() (*string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return new(id.String()), nil
}
