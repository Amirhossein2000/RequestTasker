package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"

	"github.com/Amirhossein2000/RequestTasker/internal/app/api"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/entities"
)

func (h *Handler) PostTask(ctx context.Context, request api.PostTaskRequestObject) (api.PostTaskResponseObject, error) {
	err := validatePostTaskRequestObject(request)
	if err != nil {
		return api.PostTask400JSONResponse{
			Message: err.Error(),
		}, nil
	}

	var headers map[string]string
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
		h.logger.Error("createTaskUseCase failed",
			"error", err,
		)
		return api.PostTask500Response{}, nil
	}

	return api.PostTask201JSONResponse{
		Id: publicId.String(),
	}, nil
}

func convertHeadersForRequest(headers map[string]any) map[string]string {
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

func validatePostTaskRequestObject(request api.PostTaskRequestObject) error {
	_, err := url.ParseRequestURI(request.Body.Url)
	if err != nil {
		return errors.New("url is invalid")
	}

	if !slices.Contains(methods, request.Body.Method) {
		return errors.New("method is invalid")
	}

	return nil
}

var methods = []api.HttpMethod{
	api.DELETE,
	api.GET,
	api.HEAD,
	api.OPTIONS,
	api.PATCH,
	api.POST,
	api.PUT,
}
