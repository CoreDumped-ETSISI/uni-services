package main

import (
	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
)

func (s *server) getStatus(c echo.Context) error {
	return c.JSON(200, s.status)
}

func (s *server) getHistory(c echo.Context) error {
	return c.JSON(200, s.historyCache)
}

func (s *server) route(e *echo.Echo) {
	e.Use(middleware.Logger())

	e.GET("/api/status", s.getStatus)
	e.GET("/api/history", s.getHistory)

	e.Static("/", "static")
}
