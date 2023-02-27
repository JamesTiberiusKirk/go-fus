package pages

import (
	"github.com/JamesTiberiusKirk/go-fus/examples/basic/components"
	"github.com/JamesTiberiusKirk/go-fus/fus"
	"github.com/JamesTiberiusKirk/go-fus/fusint"
	"github.com/labstack/echo/v4"
)

type HomepageData struct {
	SecretString string
	Cmp          fusint.ComponentInterface
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
		Components: []fusint.ComponentInterface{
			components.NewHomePageComponent(components.HomePageCompoentParams{}),
		},
	}

	return homepage
}

func (h *Homepage) getPageData(c echo.Context) (interface{}, error) {
	return HomepageData{
		SecretString: "",
	}, nil
}
