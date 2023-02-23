package pages

import (
	"github.com/JamesTiberiusKirk/go-fus/examples/basic/components"
	"github.com/JamesTiberiusKirk/go-fus/fus"
	"github.com/JamesTiberiusKirk/go-fus/renderer"
	"github.com/labstack/echo/v4"
)

type HomepageData struct {
	SecretString string
	Cmp          renderer.ComponentInterface
}

type Homepage struct {
	Secret string
}

func NewHomepage() *fus.Page {
	deps := &Homepage{}
	return &fus.Page{
		Deps:        deps,
		Title:       "Home page",
		ID:          "homepage",
		URI:         "/",
		Frame:       "frame",
		Template:    "homepage.gohtml",
		GetPageData: deps.GetPageData,
	}
}

func (h *Homepage) GetPageData(c echo.Context) (interface{}, error) {
	return HomepageData{
		SecretString: h.Secret,
		Cmp:          components.NewHomePageComponent(c, components.HomePageCompoentParams{}),
	}, nil
}
