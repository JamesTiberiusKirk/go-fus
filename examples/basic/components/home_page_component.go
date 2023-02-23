package components

import (
	"github.com/JamesTiberiusKirk/go-fus/fus"
	"github.com/labstack/echo/v4"
)

type HomePageCompoent struct {
	*fus.Component
	ComponentSecret string
}

type HomePageCompoentParams struct {
	Data string
}

func NewHomePageComponent(c echo.Context,
	params HomePageCompoentParams) *HomePageCompoent {
	return &HomePageCompoent{
		Component: fus.NewComponent(c,
			"homePageComponent",
			"homepage_component.gohtml",
			func() (interface{}, error) {
				return params, nil
			},
		),
		ComponentSecret: "test-secret",
	}
}
