package handlers

import (
	"context"
	"errors"

	"github.com/Amirhossein2000/RequestTasker/internal/app/api"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/common"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

func (h *Handler) GetTaskId(ctx context.Context, request api.GetTaskIdRequestObject) (api.GetTaskIdResponseObject, error) {
	publicId, err := validateGetTaskIdRequestObject(request)
	if err != nil {
		return api.GetTaskId400JSONResponse{
			Message: err.Error(),
		}, nil
	}

	task, taskStatus, taskResult, err := h.getTaskUseCase.Execute(ctx, publicId)
	if err != nil {
		switch err {
		case common.ErrNotFound:
			return api.GetTaskId404Response{}, nil

		default:
			// TODO: log
			return api.GetTaskId500Response{}, nil
		}
	}

	resp := api.GetTaskId200JSONResponse{
		Id:     task.PublicID().String(),
		Status: api.TaskStatus(taskStatus.Status()),
	}

	if taskResult != nil {
		resp.HttpStatusCode = lo.ToPtr(taskResult.StatusCode())
		resp.Length = lo.ToPtr(int(taskResult.Length()))
		resp.Headers = lo.ToPtr(convertHeadersForResponse(taskResult.Headers()))
	}

	return resp, nil
}

func convertHeadersForResponse(headers map[string]string) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range headers {
		result[key] = value
	}

	return result
}

func validateGetTaskIdRequestObject(request api.GetTaskIdRequestObject) (uuid.UUID, error) {
	if request.Id == "" {
		return uuid.Nil, errors.New("id is required")
	}

	publicId, err := uuid.Parse(request.Id)
	if err != nil {
		return uuid.Nil, errors.New("invalid id format")
	}

	return publicId, nil
}
