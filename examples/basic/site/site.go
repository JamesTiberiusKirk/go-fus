package site

import (
	"encoding/json"
	"html/template"

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

	publicPages := []fus.PageInterface{
		pages.NewHomepage(),
	}

	site.SetupTemplating("examples/basic/templates",
		publicPages,
		nil,
		pages.NewNotFoundPage(),
		map[string]string{
			Frame: "frame.gohtml",
		},
		template.FuncMap{
			"toJson": func(val interface{}) string {
				bytes, _ := json.MarshalIndent(val, "", "	")
				return string(bytes)
			},
		},
	)
	return site
}
