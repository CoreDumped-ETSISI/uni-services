package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func (s *server) getCafe(c echo.Context) error {
	menu := s.getCachedCafeMenu()

	if menu == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "El menu está desactualizado.")
	}

	return c.JSON(http.StatusOK, menu)
}

func (s *server) getTodaysMenu(c echo.Context) error {
	wd := time.Now().Weekday() - 1

	if wd < 0 {
		wd += 7
	}

	if int(wd) >= 5 {
		return echo.NewHTTPError(http.StatusBadRequest, "Hoy no hay clase.")
	}

	menu := s.getCachedCafeMenu()

	if menu == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "El menu está desactualizado.")
	}

	if !isCacheValid(menu.From, menu.To, time.Now()) {
		// Maybe invalidate cache?
		return echo.NewHTTPError(http.StatusServiceUnavailable, "El menu todavía no ha sido actualizado.")
	}

	return c.JSON(http.StatusOK, menu.Menu[wd])
}

func (s *server) getTomorrowsMenu(c echo.Context) error {
	wd := time.Now().Weekday()

	if int(wd) >= 5 {
		return echo.NewHTTPError(http.StatusBadRequest, "Mañana no hay clase.")
	}

	menu := s.getCachedCafeMenu()

	if menu == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "El menu está desactualizado.")
	}

	if !isCacheValid(menu.From, menu.To, time.Now()) {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "El menu todavía no ha sido actualizado.")
	}

	return c.JSON(http.StatusOK, menu.Menu[wd])
}

func (s *server) getStats(c echo.Context) error {
	if s.cache == nil {
		return c.String(503, "no cache")
	}

	return c.JSON(200, s.cache)
}

func route(s *server, e *echo.Echo) {
	e.GET("/all", s.getCafe)
	e.GET("/today", s.getTodaysMenu)
	e.GET("/tomorrow", s.getTomorrowsMenu)

	e.GET("/", s.getStats)

	e.Use(middleware.Logger())
}
