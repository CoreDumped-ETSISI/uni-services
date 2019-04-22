package main

import (
	"strconv"
	"sync"

	"github.com/CoreDumped-ETSISI/uni-services/emt/api"
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
	f := func(id int, o *[]api.Bus, wg *sync.WaitGroup) {
		*o, _ = s.emt.GetStopEstimates(id)
		wg.Done()
	}

	var busc []api.Bus
	var buss []api.Bus
	var busp []api.Bus

	var wg sync.WaitGroup

	wg.Add(3)

	go f(4281, &busc, &wg)
	go f(4702, &buss, &wg)
	go f(2613, &busp, &wg)

	wg.Wait()

	return c.JSON(200, UniversityStops{
		SentidoSierra: buss,
		SentidoConde:  busc,
		Puente:        busp,
	})
}

func (s *server) route(e *echo.Echo) {
	e.GET("/api/stop/:stop", s.getEstimatesForStop)
	e.GET("/api/stop", s.getEstimatesForUni)
}
