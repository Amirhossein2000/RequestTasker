package handlers

import (
	"RequestTasker/internal/app/api"
	"RequestTasker/internal/domain/common"
	"RequestTasker/internal/domain/entities"
	"context"
	"fmt"

	"github.com/samber/lo"
)

func (h *Handler) PostTask(ctx context.Context, request api.PostTaskRequestObject) (api.PostTaskResponseObject, error) {
	task := entities.NewTask( //TODO: handle bad request in this func
		request.Body.Url,
		string(request.Body.Method),
		convertHeadersForRequest(*request.Body.Headers),
		*request.Body.Body,
	)

	publicId, err := h.createTaskUseCase.Execute(ctx, task)
	if err != nil {
		switch err {
		case common.ErrInternal:
			return api.PostTask500Response{}, nil
		}
	}

	return api.PostTask201JSONResponse{
		Id: lo.ToPtr(publicId.String()),
	}, nil
}

func convertHeadersForRequest(headers map[string]interface{}) map[string]string {
	result := make(map[string]string)

	for key, value := range headers {
		if strValue, ok := value.(string); ok {
			result[key] = strValue
		} else {
			result[key] = fmt.Sprintf("%v", value)
		}
	}

	return result
}
