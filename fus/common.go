package fus

import "github.com/labstack/echo/v4"

type GetDataFunc func(c echo.Context) (interface{}, error)
