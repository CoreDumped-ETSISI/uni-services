package main

import (
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	route(e)

	e.Logger.Fatal(e.Start(":8080"))
}
