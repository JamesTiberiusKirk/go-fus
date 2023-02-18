package fus

import (
	"html/template"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/JamesTiberiusKirk/go-fus/fus/auth"
	"github.com/JamesTiberiusKirk/go-fus/renderer"
	"github.com/labstack/echo/v4"

	echoMw "github.com/labstack/echo/v4/middleware"
)

type RoutesMap map[string]string

// Site site struct with config and dependencies.
type Site struct {
	Name            string
	Dev             bool
	RootSitePath    string
	PublicPages     []*Page
	AuthedPages     []*Page
	NotFoundPage    *Page
	StaticFolders   map[string]string
	SpaSites        []*SPA
	SessionManager  *auth.AuthInterface
	Echo            *echo.Echo
	FrameTemplates  map[string]string
	TemplateFuncs   template.FuncMap
	AvailableRoutes map[string]RoutesMap
	TemplateRoot    string
}

// Serve to start the server.
func (s *Site) Serve() {
	s.buildRenderer()

	s.mapPages(&s.PublicPages)
	s.mapPages(&s.AuthedPages, sessionAuthMiddleware(s.SessionManager))

	// Mapping 404 page
	s.Echo.GET(s.RootSitePath+s.NotFoundPage.Path,
		s.NotFoundPage.GetPageHandler(http.StatusNotFound, *s.SessionManager, s.AvailableRoutes))

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

func (s *Site) buildRenderer() {
	s.Echo.Renderer = renderer.New(renderer.Config{
		Root:         s.TemplateRoot,
		Master:       s.FrameTemplates["frame"],
		NoFrame:      s.FrameTemplates["no_frame"],
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

func (s *Site) mapPages(pages *[]*Page, middlewares ...echo.MiddlewareFunc) {
	for _, p := range *pages {
		route := s.RootSitePath + p.Path
		s.AvailableRoutes[s.Name][p.MenuID] = route
	}

	for _, p := range *pages {
		route := s.RootSitePath + p.Path
		s.Echo.GET(route, p.GetPageHandler(http.StatusOK, *s.SessionManager, s.AvailableRoutes), middlewares...)

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
