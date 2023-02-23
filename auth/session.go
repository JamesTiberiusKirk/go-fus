package auth

import (
	"fmt"

	"github.com/JamesTiberiusKirk/go_web_template/models"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

// SessionManager maintains a record of open sessions.
type SessionManager struct {
	Jar             *sessions.CookieStore
	sessionName     string
	sessionLifeSpan int
}

// New returns an instantiated session manager.
func New(sessionName string, sessionLifeSpan int) *SessionManager {
	authKey := securecookie.GenerateRandomKey(64)
	encryptionKey := securecookie.GenerateRandomKey(32)

	return &SessionManager{
		Jar:             sessions.NewCookieStore(authKey, encryptionKey),
		sessionName:     sessionName,
		sessionLifeSpan: sessionLifeSpan,
	}
}

// InitSession will store a new session or refresh an existing one.
func (m *SessionManager) InitSession(user interface{}, c echo.Context) error {
	sess, err := m.Jar.Get(c.Request(), m.sessionName)
	if err != nil {
		return fmt.Errorf("error getting session: %w", err)
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   m.sessionLifeSpan,
		HttpOnly: true,
	}

	sess.Values["user"] = user

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return fmt.Errorf("error saving session: %w", err)
	}

	return nil
}

// TerminateSession will cease tracking the session for the current user.
func (m *SessionManager) TerminateSession(c echo.Context) error {
	sess, err := m.Jar.Get(c.Request(), m.sessionName)
	if err != nil {
		return fmt.Errorf("error getting session: %w", err)
	}

	// MaxAge < 0 means delete imediately
	sess.Options.MaxAge = -1
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return fmt.Errorf("error saving session: %w", err)
	}

	return nil
}

// IsAuthenticated checks that a provided request is born from an active session.
// As long as there is an active session, true is returned, else false.
func (m *SessionManager) IsAuthenticated(c echo.Context) (bool, error) {
	sess, err := m.Jar.Get(c.Request(), m.sessionName)
	if err != nil {
		return false, fmt.Errorf("error getting session: %w", err)
	}

	return sess.Values["email"] != nil, nil
}

// GetUser checks that a provided request is born from an active session.
// As long as there is an active session, User is returned, else empty User.
func (m *SessionManager) GetUser(c echo.Context) (interface{}, error) {
	sess, err := m.Jar.Get(c.Request(), m.sessionName)
	if err != nil {
		return models.User{}, err
	}

	if sess.Values == nil {
		return models.User{}, nil
	}

	user, ok := sess.Values["user"]
	if !ok {
		return models.User{}, nil
	}

	return user, nil
}
