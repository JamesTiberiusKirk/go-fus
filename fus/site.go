package fus

import (
	"html/template"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/JamesTiberiusKirk/go-fus/auth"
	"github.com/JamesTiberiusKirk/go-fus/renderer"
	"github.com/labstack/echo/v4"

	echoMw "github.com/labstack/echo/v4/middleware"
)

type RoutesMap map[string]string

// Site site struct with config and dependencies.
type Site struct {
	// Echo - Instance of echo
	Echo *echo.Echo

	// Name - Name of the site, will be used in AvailableRoutes.
	Name string

	// Dev - Dev mode, used to toggle proxying to local SPA dev servers
	Dev bool

	// RootSitePath - The root path of the site
	RootSitePath string

	// PublicPages - Public pages
	PublicPages []*Page

	// AuthedPages - Pages to be authenticated with the session manager
	AuthedPages []*Page

	// NotFoundPage - 404 page
	NotFoundPage *Page

	// StaticFolders - Static folders to serve
	StaticFolders map[string]string

	// SpaSites - single page apps to initiate
	SpaSites []*SPA

	// SessionManager - Session manager to use
	SessionManager auth.SessionInterface

	// FrameTemplates - List of frames which could be used.
	// For having a no frame option, just create an empty frame.
	FrameTemplates map[string]string

	// TemplateFuncs - is for defining any extra template functions
	TemplateFuncs template.FuncMap

	// AvailableRoutes - this is for defining a map of availabe routes which exist outside
	//	the site, which then would be made available in the template
	// E.G. defining available routes in a json API then being able to access it in template
	//	with `{{ routes.api.helloWorldRoute }}`
	AvailableRoutes map[string]RoutesMap

	// TemplateRoot - root folder where the templates are located
	TemplateRoot string
}

func NewSite(e *echo.Echo, name string, dev bool, rootSitePath string,
	availableRoutes map[string]RoutesMap) Site {
	if availableRoutes == nil {
		availableRoutes = map[string]RoutesMap{}
	}
	availableRoutes[name] = RoutesMap{}

	return Site{
		Echo:            e,
		Name:            name,
		Dev:             dev,
		RootSitePath:    rootSitePath,
		AvailableRoutes: availableRoutes,
	}
}

// SetupTemplating Setting up templating option for the site.
func (s *Site) SetupTemplating(
	templateRoot string,
	publicPages, authedPages []*Page,
	notFoundPage *Page,
	frameTemplates map[string]string,
	templateFuncs template.FuncMap) {
	//
	// if frameTemplates == nil {
	// 	frameTemplates = map[string]string{}
	// }
	//
	// frameTemplates[]
	//
	s.TemplateRoot = templateRoot
	s.PublicPages = publicPages
	s.AuthedPages = authedPages
	s.NotFoundPage = notFoundPage
	s.FrameTemplates = frameTemplates
	s.TemplateFuncs = templateFuncs
}

// SetupSPA Mapping single page app for the site.
func (s *Site) SetupSPA(spa []*SPA) {
	s.SpaSites = spa
}

// SetupStatic Mapping static folders for the site.
func (s *Site) SetupStatic(staticFolders map[string]string) {
	s.StaticFolders = staticFolders
}

// Serve to start the server.
func (s *Site) Serve() {
	s.buildRenderer()

	s.MapPages(&s.PublicPages)
	// s.MapPages(&s.AuthedPages, sessionAuthMiddleware(s.SessionManager))

	// Mapping 404 page
	s.Echo.GET(s.RootSitePath+s.NotFoundPage.URI,
		s.NotFoundPage.getPageHandler(http.StatusNotFound, s.SessionManager, s.AvailableRoutes))

	s.mapStatic()
	s.mapSPA()
}

// GetRoutes to get routes which have been made in the server.
func (s *Site) GetRoutes() RoutesMap {
	return s.AvailableRoutes[s.Name]
}

// SetRoutes which would be used in the templating engine.
func (s *Site) SetRoutes(t string, r RoutesMap) {
	s.AvailableRoutes[t] = r
}

// MapPages - takes a pointer to a list of Pages and any middlewares to be used when initiating them.
func (s *Site) MapPages(pages *[]*Page, middlewares ...echo.MiddlewareFunc) {
	if pages == nil {
		return
	}

	if s.AvailableRoutes != nil {
		for _, p := range *pages {
			route := s.RootSitePath + p.URI
			s.AvailableRoutes[s.Name][p.ID] = route
		}
	}

	for _, p := range *pages {
		route := s.RootSitePath + p.URI
		pageHandler := p.getPageHandler(http.StatusOK, s.SessionManager, s.AvailableRoutes)
		s.Echo.GET(route, pageHandler, middlewares...)

		if p.PostHandler != nil {
			s.Echo.POST(route, p.PostHandler, middlewares...)
		}

		if p.DeleteHandler != nil {
			s.Echo.DELETE(route, p.DeleteHandler, middlewares...)
		}

		if p.PutHandler != nil {
			s.Echo.PUT(route, p.PutHandler, middlewares...)
		}
	}
}

func (s *Site) buildRenderer() {
	s.Echo.Renderer = renderer.New(renderer.Config{
		Root:         s.TemplateRoot,
		Frames:       s.FrameTemplates,
		Funcs:        s.TemplateFuncs,
		DisableCache: true,
	})
}

func (s *Site) mapSPA(_ ...echo.MiddlewareFunc) {
	for _, spa := range s.SpaSites {
		route := s.RootSitePath + spa.Path

		switch s.Dev {
		case true:
			proxy := httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: "http",
				Host:   spa.Dev.Host,
			})
			s.Echo.Any(spa.Path+"*", echo.WrapHandler(proxy))
		case false:
			group := s.Echo.Group(route)
			group.Use(echoMw.StaticWithConfig(echoMw.StaticConfig{
				Root:   spa.Dist,
				Index:  spa.Index,
				Browse: spa.Routing,
				HTML5:  true,
			}))
		}

		s.AvailableRoutes[s.Name][spa.MenuID] = route
	}
}

func (s *Site) mapStatic() {
	for k, v := range s.StaticFolders {
		s.Echo.Static(k, v)
	}
}
