package main

import (
	"github.com/JamesTiberiusKirk/go-fus/examples/basic/site"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())

	siteServer := site.NewTestSie(e, "Secret string")

	siteServer.Serve()

	e.Logger.Fatal(e.Start(":3000"))
}
