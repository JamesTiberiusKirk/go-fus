package pages

import (
	"github.com/JamesTiberiusKirk/go-fus/fus"
	"github.com/labstack/echo/v4"
)

func NewNotFoundPage() *fus.Page {
	return &fus.Page{
		ID:       "notFoundPage",
		Title:    "NotFound",
		Frame:    "frame",
		URI:      "/*",
		Template: "not_found.gohtml",
		Deps:     nil,
		GetPageData: func(c echo.Context) (interface{}, error) {
			return nil, nil
		},
	}
}
