package site

import (
	"github.com/JamesTiberiusKirk/go-fus/examples/basic/pages"
	"github.com/JamesTiberiusKirk/go-fus/fus"
	"github.com/labstack/echo/v4"
)

const (
	Frame = "frame"
)

type TestSite struct {
	fus.Site
	someSecret string
}

func NewTestSie(echo *echo.Echo, secret string) *TestSite {
	site := &TestSite{
		Site:       fus.NewSite(echo, "testSite", true, "", map[string]fus.RoutesMap{}),
		someSecret: secret,
	}

	site.SetupTemplating("examples/basic/templates",
		[]*fus.Page{
			pages.NewHomepage(),
		},
		nil,
		pages.NewNotFoundPage(),
		map[string]string{
			Frame: "frame.gohtml",
		},
		nil,
	)
	return site
}
