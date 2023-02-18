package fus

import (
	"errors"

	"github.com/JamesTiberiusKirk/go-fus/auth"
	"github.com/labstack/echo/v4"
)

// MetaData is used to give certain page meta data and basic params to each template.
type MetaData struct {
	MenuID   string
	Title    string
	URLError string
	Success  string
}

// Page is used by every page in a site
// Deps being each page's own struct for dependencies, might not even be needed.
type Page struct {
	MenuID        string
	Title         string
	Frame         bool
	Path          string
	Template      string
	Deps          interface{}
	GetPageData   func(c echo.Context) any
	PostHandler   echo.HandlerFunc
	DeleteHandler echo.HandlerFunc
	PutHandler    echo.HandlerFunc
}

const (
	UseFrameName = "frame"
)

// GetPageHandler is a get handler which uses the echo Render function.
func (p *Page) GetPageHandler(httpStatus int, session auth.SessionInterface,
	routesMap map[string]RoutesMap) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set(UseFrameName, p.Frame)
		var auth echo.Map

		user, err := session.GetUser(c)
		if err != nil {
			if errors.Is(err, errors.New("securecookie: the value is not valid")) {
				return err
			}

			auth = echo.Map{}
		} else {
			auth = echo.Map{
				"email":    user.Email,
				"username": user.Username,
			}
		}

		echoData := echo.Map{
			"meta":   p.buildBasePageMetaData(c),
			"auth":   auth,
			"routes": routesMap,
		}

		pageData := p.GetPageData(c)
		if pageData != nil {
			echoData["data"] = pageData
		}

		err = c.Render(httpStatus, p.Template, echoData)
		if err != nil {
			return err
		}

		return nil
	}
}

func (p *Page) buildBasePageMetaData(c echo.Context) MetaData {
	return MetaData{
		MenuID:   p.MenuID,
		Title:    p.Title,
		URLError: c.QueryParam("error"),
		Success:  c.QueryParam("success"),
	}
}
