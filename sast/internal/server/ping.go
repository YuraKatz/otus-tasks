package server

import (
	"net/http"

	"github.com/labstack/echo"
)

func (s *Server) pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
