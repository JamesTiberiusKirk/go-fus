package fus

import (
	"net/http"

	"github.com/JamesTiberiusKirk/go-fus/fus/auth"
	echoSession "github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func sessionAuthMiddleware(auth auth.AuthInterface) echo.MiddlewareFunc {
	return echoSession.MiddlewareWithConfig(echoSession.Config{
		Skipper: func(c echo.Context) bool {
			skip := auth.IsAuthenticated(c)
			if !skip {
				_ = c.Redirect(http.StatusSeeOther, "/login")
			}

			return skip
		},
		Store: auth.GetJar(),
	})
}
