package server

import (
	"context"
	"fmt"
	"sasthw/internal/config"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Server struct {
	ctx  context.Context
	cfg  config.Config
	echo *echo.Echo
}

func New(cfg config.Config) (*Server, error) {
	return &Server{
		ctx:  context.Background(),
		cfg:  cfg,
		echo: echo.New(),
	}, nil
}

func (s *Server) onStart() {
	s.echo.HideBanner = true
	s.echo.Debug = s.cfg.Debug
}

func (s *Server) ListenAndServe() error {

	s.onStart()

	s.echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	s.echo.File("/", "tree.html")
	s.echo.GET("/ping", s.pingHandler)

	s.echo.PUT("/tree", s.treeHandler)

	static := s.echo.Group("/s")
	{
		static.Static("/", s.cfg.StaticPath)
	}

	return s.echo.Start(fmt.Sprintf(":%d", s.cfg.HTTPPort))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
