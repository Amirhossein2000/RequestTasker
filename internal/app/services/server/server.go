package server

import (
	"context"

	"github.com/Amirhossein2000/RequestTasker/internal/app/api"

	"github.com/labstack/echo/v4"
)

type Server struct {
	addr        string
	handlers    api.StrictServerInterface
	middlewares []api.StrictMiddlewareFunc
	e           *echo.Echo
}

func NewServer(
	addr string,
	handlers api.StrictServerInterface,
	Middlewares []api.StrictMiddlewareFunc,
) *Server {
	return &Server{
		addr:        addr,
		handlers:    handlers,
		middlewares: Middlewares,
	}
}

func (s *Server) Start() error {
	e := echo.New()
	handlers := api.NewStrictHandler(s.handlers, s.middlewares)
	api.RegisterHandlers(e, handlers)
	s.e = e
	return s.e.Start(s.addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}
