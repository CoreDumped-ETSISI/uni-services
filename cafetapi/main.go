package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/robfig/cron"
)

const (
	// CacheTimeout is the cache invalidation timeout in hours.
	CacheTimeout = 4
)

type server struct {
	tableServer string
	cache       *menuCache
}

func main() {
	s := &server{}

	s.tableServer = os.Getenv("PDF_TABLE_SERVER")
	if s.tableServer == "" {
		s.tableServer = "localhost"
	}

	e := echo.New()

	// Keep cache warm
	c := cron.New()
	updateFunc := s.updateCache(e.Logger)
	c.AddFunc(fmt.Sprintf("0 0 */%v * * 2-5", CacheTimeout/2), updateFunc) // General, 2hr, keep cache warm
	c.AddFunc("0 */15 * * * 1", updateFunc)                                // Try to get latest one on monday, every 15 min

	c.Start()

	updateFunc()

	route(s, e)
	e.Logger.Fatal(e.Start(":8889"))
}
