package main

import (
	"strconv"
	"time"

	"github.com/labstack/echo/middleware"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
)

func (s *server) getStatus(c echo.Context) error {
	return c.JSON(200, s.status)
}

func (s *server) getHistory(c echo.Context) error {
	limit := c.QueryParam("limit")
	last := c.QueryParam("last")

	sesh := s.mongo.Clone()
	defer sesh.Close()

	col := sesh.DB("etsisi-telegram-bot").C("status_history")

	var his []serviceHistory

	var q *mgo.Query

	if last != "" {
		dur, _ := time.ParseDuration(last)
		q = col.Find(bson.M{
			"_id": bson.M{
				"$gte": time.Now().Add(-dur),
			},
		})
	} else {
		q = col.Find(nil)
	}

	if limit != "" {
		l, _ := strconv.Atoi(limit)
		q = q.Limit(l)
	}

	err := q.All(&his)

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
