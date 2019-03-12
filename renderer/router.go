package main

import (
	"io"

	"github.com/CoreDumped-ETSISI/uni-services/renderer/horario"
	"github.com/CoreDumped-ETSISI/uni-services/renderer/salas"

	"github.com/labstack/echo"
)

type renderer func(body io.Reader, response io.Writer) error

func renderSalas(c echo.Context) error {
	defer c.Request().Body.Close()
	c.Response().WriteHeader(200)
	c.Response().Header().Set("Content-Type", "image/png")

	return salas.RenderSalas(c.Request().Body, c.Response())
}

func renderHorario(c echo.Context) error {
	class := c.Param("group")

	defer c.Request().Body.Close()
	c.Response().WriteHeader(200)
	c.Response().Header().Set("Content-Type", "image/png")

	return horario.RenderImage(class, c.Request().Body, c.Response())
}

func route(e *echo.Echo) {
	e.POST("/api/salas", renderSalas)
	e.POST("/api/horario", renderHorario)
}
