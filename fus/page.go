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

// PageInterface interface which is implemented by the Page struct bellow
// TODO: write docs for the rest of the funcs
type PageInterface interface {
	GetID() string
	GetTitle() string
	GetFrame() string
	GetURI() string
	GetPageDataHandler() func(c echo.Context) (interface{}, error)
	GetPostHandler() echo.HandlerFunc
	GetDeleteHandler() echo.HandlerFunc
	GetPutHandler() echo.HandlerFunc

	// This one is for internal use
	GetPageHandler(httpStatus int, session auth.SessionInterface, routesMap map[string]RoutesMap) echo.HandlerFunc
	GetComponents() []ComponentInterface
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

	PageDataHandler func(c echo.Context) (interface{}, error)
	PostHandler     echo.HandlerFunc
	DeleteHandler   echo.HandlerFunc
	PutHandler      echo.HandlerFunc

	Components []ComponentInterface
}

const (
	UseFrameName = "frame"
)

// GetPageHandler is a get handler which uses the echo Render function.
func (p *Page) GetPageHandler(httpStatus int, session auth.SessionInterface,
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

		handler := p.GetPageDataHandler()
		pageData, err := handler(c)

		echoData["data"] = pageData
		echoData["error"] = err

		compoentsMap := echo.Map{}
		compoennts := p.GetComponents()
		for _, component := range compoennts {
			component.SetContext(c)
			compoentsMap[component.GetID()] = component
		}

		echoData["c"] = compoentsMap

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

func (p *Page) GetID() string                       { return p.ID }
func (p *Page) GetTitle() string                    { return p.Title }
func (p *Page) GetFrame() string                    { return p.Frame }
func (p *Page) GetURI() string                      { return p.URI }
func (p *Page) GetPostHandler() echo.HandlerFunc    { return p.PostHandler }
func (p *Page) GetDeleteHandler() echo.HandlerFunc  { return p.DeleteHandler }
func (p *Page) GetPutHandler() echo.HandlerFunc     { return p.PutHandler }
func (p *Page) GetComponents() []ComponentInterface { return p.Components }

func (p *Page) GetPageDataHandler() func(c echo.Context) (interface{}, error) {
	return p.PageDataHandler
}
