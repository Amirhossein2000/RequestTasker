package handlers

import (
	"context"
	"fmt"

	"github.com/Amirhossein2000/RequestTasker/internal/app/api"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/entities"
)

func (h *Handler) PostTask(ctx context.Context, request api.PostTaskRequestObject) (api.PostTaskResponseObject, error) {
	headers := make(map[string]string)
	body := ""

	if request.Body.Headers != nil {
		headers = convertHeadersForRequest(*request.Body.Headers)
	}

	if request.Body.Body != nil {
		body = *request.Body.Body
	}

	task := entities.NewTask(
		request.Body.Url,
		string(request.Body.Method),
		headers,
		body,
	)

	publicId, err := h.createTaskUseCase.Execute(ctx, task)
	if err != nil {
		return api.PostTask500Response{}, nil
	}

	return api.PostTask201JSONResponse{
		Id: publicId.String(),
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
