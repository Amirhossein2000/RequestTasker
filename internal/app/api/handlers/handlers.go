package handlers

import (
	"github.com/Amirhossein2000/RequestTasker/internal/app/services/logger"
	"github.com/Amirhossein2000/RequestTasker/internal/app/usecases"
)

type Handler struct {
	logger            *logger.Logger
	apiKey            string
	createTaskUseCase usecases.CreateTaskUseCase
	getTaskUseCase    usecases.GetTaskUseCase
}

func NewHandler(
	logger *logger.Logger,
	apiKey string,
	createTaskUseCase usecases.CreateTaskUseCase,
	getTaskUseCase usecases.GetTaskUseCase,
) *Handler {
	return &Handler{
		logger:            logger,
		apiKey:            apiKey,
		createTaskUseCase: createTaskUseCase,
		getTaskUseCase:    getTaskUseCase,
	}
}
