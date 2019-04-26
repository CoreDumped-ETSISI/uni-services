package main

import (
	"strconv"
	"time"

	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
)

func (s *server) getStatus(c echo.Context) error {
	return c.JSON(200, s.status)
}

func (s *server) getHistory(c echo.Context) error {
	limit := c.QueryParam("limit")
	last := c.QueryParam("last")

	var his []serviceHistory

	q := s.postgres.Model(&his)

	if last != "" {
		dur, _ := time.ParseDuration(last)
		t := time.Now().Add(-dur)

		q = q.Where("timestamp > ?", t)
	}

	if limit != "" {
		l, _ := strconv.Atoi(limit)
		q = q.Limit(l)
	}

	err := q.Select(&his)

	if err != nil {
		return err
	}

	var thinhis []*serviceHistory
	days := map[string]int{}

	for i := range his {
		if days[his[i].URL] != his[i].Timestamp.Day() || !his[i].Up {
			days[his[i].URL] = his[i].Timestamp.Day()
			thinhis = append(thinhis, &his[i])
		}
	}

	return c.JSON(200, thinhis)
}

func (s *server) route(e *echo.Echo) {
	e.Use(middleware.Logger())

	e.GET("/api/status", s.getStatus)
	e.GET("/api/history", s.getHistory)

	e.Static("/", "static")
}
