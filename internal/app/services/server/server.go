// server.go
package server

import (
	"context"

	"github.com/Amirhossein2000/RequestTasker/internal/app/api"

	"github.com/labstack/echo/v4"
)

type Server struct {
	Addr        string
	Handlers    api.StrictServerInterface
	Middlewares []api.StrictMiddlewareFunc
	e           *echo.Echo
}

func NewServer(
	addr string,
	handlers api.StrictServerInterface,
	Middlewares []api.StrictMiddlewareFunc,
) *Server {
	return &Server{
		Addr:        addr,
		Handlers:    handlers,
		Middlewares: Middlewares,
	}
}

func (s *Server) Start() error {
	e := echo.New()
	handlers := api.NewStrictHandler(s.Handlers, s.Middlewares)
	api.RegisterHandlers(e, handlers)
	s.e = e
	return s.e.Start(s.Addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}
