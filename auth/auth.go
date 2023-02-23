package auth

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type SessionInterface interface {
	InitSession(user interface{}, c echo.Context) error
	TerminateSession(c echo.Context) error
	IsAuthenticated(c echo.Context) bool
	GetUser(c echo.Context) (interface{}, error)
	GetJar() *sessions.CookieStore
}
