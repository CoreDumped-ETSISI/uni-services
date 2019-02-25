package main

import (
	"github.com/labstack/echo"
	"github.com/robfig/cron"
)

func main() {
	e := echo.New()
	server := New()

	server.log = e.Logger

	server.launchPrecheck()

	server.route(e)
	c := cron.New()

	c.AddFunc("0 */1 * * * *", server.checkAllServices)

	c.Start()

	e.Logger.Fatal(e.Start(":8889"))
}
