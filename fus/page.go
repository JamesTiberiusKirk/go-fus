package fus

import (
	"errors"

	"github.com/JamesTiberiusKirk/go-fus/auth"
	"github.com/labstack/echo/v4"
)

// MetaData is used to give certain page meta data and basic params to each template.
type MetaData struct {
	ID       string
	Title    string
	URLError string
	Success  string
}

// Page is used by every page in a site
// Deps being each page's own struct for dependencies, might not even be needed.
type Page struct {
	// ID - used for routes mapping, can be used for a menu in the frame.
	// Accessible in page meta data.
	ID string

	// Title - semantic title of the page
	// Accessible in page meta data.
	Title string

	// Frame - which frame to use for the page
	Frame string

	// URI - on what the page is served
	URI string

	// Template file to be used. Needs to be inside the template directory in Site.
	Template string

	// Deps - any dependencies the page needs to use.
	Deps interface{}

	// GetPageData function to get any data used in the teamplate.
	// Both of the returns are passed to the template.
	GetPageData func(c echo.Context) (interface{}, error)

	// PostHandler a POST handler for the page.
	// Can be nil to omit the definition of one.
	PostHandler echo.HandlerFunc

	// DeleteHandler a DELETE handler for the page.
	// Can be nil to omit the definition of one.
	DeleteHandler echo.HandlerFunc

	// PutHandler a PUT handler for the page.
	// Can be nil to omit the definition of one.
	PutHandler echo.HandlerFunc
}

const (
	UseFrameName = "frame"
)

// getPageHandler is a get handler which uses the echo Render function.
func (p *Page) getPageHandler(httpStatus int, session auth.SessionInterface,
	routesMap map[string]RoutesMap) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set(UseFrameName, p.Frame)

		echoData := echo.Map{
			"meta":   p.buildBasePageMetaData(c),
			"routes": routesMap,
		}

		if session != nil {
			user, err := session.GetUser(c)
			if err != nil {
				if errors.Is(err, errors.New("securecookie: the value is not valid")) {
					return err
				}
			}
			echoData["user"] = user
		}

		pageData, err := p.GetPageData(c)
		echoData["data"] = pageData
		echoData["error"] = err

		err = c.Render(httpStatus, p.Template, echoData)
		if err != nil {
			return err
		}

		return nil
	}
}

func (p *Page) buildBasePageMetaData(c echo.Context) MetaData {
	return MetaData{
		ID:       p.ID,
		Title:    p.Title,
		URLError: c.QueryParam("error"),
		Success:  c.QueryParam("success"),
	}
}
