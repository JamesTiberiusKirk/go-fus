package fus

import (
	"net/http"

	"github.com/JamesTiberiusKirk/go-fus/auth"
	echoSession "github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// SessionAuthMiddleware middleware to use the session auth.
func SessionAuthMiddleware(auth auth.SessionInterface, loginPageURI string) echo.MiddlewareFunc {
	return echoSession.MiddlewareWithConfig(echoSession.Config{
		Skipper: func(c echo.Context) bool {
			skip := auth.IsAuthenticated(c)
			if !skip {
				_ = c.Redirect(http.StatusSeeOther, loginPageURI)
			}

			return skip
		},
		Store: auth.GetJar(),
	})
}
