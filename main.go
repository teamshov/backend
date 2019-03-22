package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORS())

	e.Static("/", "dist")
	e.File("/", "dist/index.html")

	InitDBService(e)
	InitRedisService(e)
	InitPathfindingService(e)

	e.Logger.Fatal(e.Start(":62027"))
}
