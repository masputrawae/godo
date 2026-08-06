package cookie

import (
	"context"
	"errors"
	"godo/internal/model"
	"godo/internal/repo"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSessionIsExpires = errors.New("session is expires")
)

type Cookie struct {
	Name    string
	TTL     time.Duration
	session *repo.Session
}

func New(session *repo.Session) *Cookie {
	return &Cookie{
		Name:    "session",
		TTL:     2 * time.Minute,
		session: session,
	}
}

func (c *Cookie) Set(ctx context.Context, w http.ResponseWriter, userID int) (string, string, error) {
	sessionID, err := c.RandUUID()
	if err != nil {
		return "", "", err
	}
	csrfToken, err := c.RandUUID()
	if err != nil {
		return "", "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     c.Name,
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(c.TTL),
		MaxAge:   int((c.TTL).Seconds()),
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	})

	err = c.session.Create(ctx, model.Session{
		ID:        sessionID,
		CSRFToken: csrfToken,
		UserID:    userID,
		ExpiresAt: time.Now().Add(c.TTL),
	})

	return sessionID, csrfToken, err
}

func (c *Cookie) Get(ctx context.Context, w http.ResponseWriter, sessionID string) (model.Session, error) {
	session, err := c.session.Find(ctx, sessionID)
	if err != nil {
		return session, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		if err := c.Remove(ctx, w, sessionID); err != nil {
			return session, err
		}
	}

	return session, nil
}

func (c *Cookie) Remove(ctx context.Context, w http.ResponseWriter, sessionID string) error {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	})

	err := c.session.Remove(ctx, sessionID)
	if err != nil && !errors.Is(err, repo.ErrSessionNotExist) {
		return err
	}
	return nil
}

func (c *Cookie) ParseCtx(r *http.Request) (model.Session, bool) {
	session, ok := r.Context().Value("session").(model.Session)
	return session, ok
}

// 2. random uuid
func (c *Cookie) RandUUID() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
