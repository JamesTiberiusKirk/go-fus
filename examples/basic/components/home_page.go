package components

import (
	"errors"

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

func NewHomePageComponent() *HomePageCompoent {
	return &HomePageCompoent{
		Component: fus.NewComponent(
			"homePageComponent",
			"homepage_component.gohtml",
			func(c echo.Context, params interface{}) (interface{}, error) {
				homePageCompoentParams, ok := params.(HomePageCompoentParams)
				if !ok {
					return nil, errors.New("missing params")
				}

				return homePageCompoentParams.Data, nil
			},
			NewListItem(),
		),
		ComponentSecret: "test-secret",
	}
}
