// server.go
package server

import (
	"RequestTasker/internal/app/api"

	"github.com/labstack/echo/v4"
)

type Server struct {
	Addr     string
	Handlers api.StrictServerInterface
}

func NewServer(addr string, handlers api.StrictServerInterface) *Server {
	return &Server{
		Addr:     addr,
		Handlers: handlers,
	}
}

func (s *Server) Start() error {
	// TODO: Auth middleware
	e := echo.New()
	handlers := api.NewStrictHandler(s.Handlers, nil)
	api.RegisterHandlers(e, handlers)
	return e.Start(s.Addr)
}
