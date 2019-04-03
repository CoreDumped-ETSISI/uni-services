package main

import (
	"strconv"

	"github.com/labstack/echo"
)

// GET /api/stop/{stop}
func (s *server) getEstimatesForStop(c echo.Context) error {
	stopval := c.Param("stop")

	sid, err := strconv.Atoi(stopval)

	if err != nil {
		return err
	}

	bus, err := s.emt.GetStopEstimates(sid)

	if err != nil {
		return err
	}

	return c.JSON(200, bus)
}

// GET /api/stop/
func (s *server) getEstimatesForUni(c echo.Context) error {
	condeid := 4281
	sierraid := 4702

	busc, errs := s.emt.GetStopEstimates(condeid)
	buss, errc := s.emt.GetStopEstimates(sierraid)

	if errs != nil && errc != nil {
		return errs
	}

	return c.JSON(200, UniversityStops{
		SentidoSierra: buss,
		SentidoConde:  busc,
	})
}

func (s *server) route(e *echo.Echo) {
	e.GET("/api/stop/:stop", s.getEstimatesForStop)
	e.GET("/api/stop", s.getEstimatesForUni)
}
