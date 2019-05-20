package main

import (
	"os"

	"github.com/labstack/echo/middleware"

	"github.com/CoreDumped-ETSISI/uni-services/emt/api"
	"github.com/labstack/echo"
)

type server struct {
	emt *api.EMT
}

func main() {
	e := echo.New()

	s := &server{}
	s.emt = api.New(os.Getenv("EMT_IDCLIENT"), os.Getenv("EMT_PASSKEY"))

	s.route(e)

	e.Use(middleware.Logger())

	e.Logger.Fatal(e.Start(":8080"))
}
