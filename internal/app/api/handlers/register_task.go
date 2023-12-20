package handlers

import (
	"RequestTasker/internal/app/api"
	"RequestTasker/internal/domain/entities"
	"context"
)

func (h *Handler) PostTask(ctx context.Context, request api.PostTaskRequestObject) (api.PostTaskResponseObject, error) {
	entities.NewTask(
		request.Body.Url,
		string(request.Body.Method),
		nil,
		"",
	)

	panic("Implement me")

}
