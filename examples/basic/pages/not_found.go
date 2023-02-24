package pages

import (
	"github.com/JamesTiberiusKirk/go-fus/fus"
	"github.com/labstack/echo/v4"
)

type NotFoundPage struct {
	*fus.Page
}

func NewNotFoundPage() *NotFoundPage {
	return &NotFoundPage{
		&fus.Page{
			ID:       "notFoundPage",
			Title:    "NotFound",
			Frame:    "frame",
			URI:      "/*",
			Template: "not_found.gohtml",
			PageDataHandler: func(c echo.Context) (interface{}, error) {
				return "", nil
			},
		},
	}
}
