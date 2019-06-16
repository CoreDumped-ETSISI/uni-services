package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	s := New()

	s.route(e)

	e.Logger.Fatal(e.Start(":8889"))
}
