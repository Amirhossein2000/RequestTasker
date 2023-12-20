// server.go
package server

import (
	"RequestTasker/internal/app/api"

	"github.com/labstack/echo/v4"
)

type Server struct {
	Addr        string
	Handlers    api.StrictServerInterface
	Middlewares []api.StrictMiddlewareFunc
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
	return e.Start(s.Addr)
}
