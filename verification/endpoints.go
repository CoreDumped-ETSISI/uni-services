package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *server) route(e *echo.Echo) {
	e.GET("/sendoff", s.sendoff)
	e.GET("/redirect", s.welcomeback)
}

// GET /sendoff
func (s *server) sendoff(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, "https://api.upm.es/todo") // TODO
}

// GET /redirect
func (s *server) welcomeback(c echo.Context) error {
	// TODO: Verify stuff

	token := c.QueryParam("token")
	const prefix = "VERIFY_SESS_"

	useridraw, err := s.redis.Get(prefix + token).Result()

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Session expired")
	}

	userid, err := strconv.Atoi(useridraw)

	if err != nil {
		return err
	}

	event := userVerified{
		UserID: userid,
	}

	data, _ := json.Marshal(event)

	err = s.redis.Publish("USER_VERIFIED", string(data)).Err()

	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Something happened")
	}

	// Delete the old token
	_ = s.redis.Del(prefix + token).Err()

	// TODO: Mejorar esto
	return c.String(http.StatusOK, "Ya puedes cerrar esta ventana.")
}
