package components

import (
	"github.com/JamesTiberiusKirk/go-fus/fus"
	"github.com/labstack/echo/v4"
)

type ListItem struct {
	*fus.Component
}

type User struct {
	Email    string
	Username string
}

func NewListItem() *ListItem {
	return &ListItem{
		fus.NewComponent(
			"listItem",
			"list_item_compoent.gohtml",
			func(c echo.Context, params interface{}) (interface{}, error) {
				user, ok := params.(User)
				if !ok {
					return User{}, nil
				}

				return user, nil
			},
		),
	}
}
