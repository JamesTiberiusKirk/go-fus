package fusint

import "github.com/labstack/echo/v4"

type ComponentInterface interface {
	GetID() string
	GetTemplate() string
	SetContext(c echo.Context)
	GenerateComponentData(parentData interface{}, params ...interface{}) (echo.Map, error)

	// getCompoents() map[string]ComponentInterface
}
