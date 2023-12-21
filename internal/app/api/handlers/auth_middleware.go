package handlers

import (
	"net/http"

	"github.com/Amirhossein2000/RequestTasker/internal/app/api"

	"github.com/labstack/echo/v4"
)

func (h *Handler) AuthMiddleware(f api.StrictHandlerFunc, operationID string) api.StrictHandlerFunc {
	return func(ctx echo.Context, request interface{}) (response interface{}, err error) {
		if h.apiKey != ctx.Request().Header.Get("Authorization") {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}
		return f(ctx, request)
	}
}
