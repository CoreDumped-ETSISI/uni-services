package main

import (
	"github.com/labstack/echo"
)

func (s *server) getStatus(c echo.Context) error {
	return c.JSON(200, s.status)
}

func (s *server) route(e *echo.Echo) {
	e.GET("/", s.getStatus)
}
