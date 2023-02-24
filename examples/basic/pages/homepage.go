package pages

import (
	"github.com/JamesTiberiusKirk/go-fus/examples/basic/components"
	"github.com/JamesTiberiusKirk/go-fus/fus"
	"github.com/labstack/echo/v4"
)

type HomepageData struct {
	SecretString string
	Cmp          fus.ComponentInterface
}

type Homepage struct {
	*fus.Page
	Secret string
	// Any dependencies here...
}

func NewHomepage() *Homepage {
	homepage := &Homepage{
		Secret: "secret string",
	}
	homepage.Page = &fus.Page{
		Title:           "Home page",
		ID:              "homepage",
		URI:             "/",
		Frame:           "frame",
		Template:        "homepage.gohtml",
		PageDataHandler: homepage.getPageData,
		// Components: map[string]fus.ComponentInterface{
		// 	"homePageComponent": components.NewHomePageComponent(c, components.HomePageCompoentParams{}),
		// },
	}

	return homepage
}

func (h *Homepage) getPageData(c echo.Context) (interface{}, error) {
	return HomepageData{
		SecretString: "",
		Cmp:          components.NewHomePageComponent(c, components.HomePageCompoentParams{}),
	}, nil
}
