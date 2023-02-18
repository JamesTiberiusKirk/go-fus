package auth

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type SessionInterface[userType interface{}] interface {
	IsAuthenticated(c echo.Context) bool
	GetUser(c echo.Context) (userType, error)
	GetJar() *sessions.CookieStore
}
