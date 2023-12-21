package handlers

import "github.com/Amirhossein2000/RequestTasker/internal/app/usecases"

type Handler struct {
	apiKey            string
	createTaskUseCase usecases.CreateTaskUseCase
	getTaskUseCase    usecases.GetTaskUseCase
}

func NewHandler(
	apiKey string,
	createTaskUseCase usecases.CreateTaskUseCase,
	getTaskUseCase usecases.GetTaskUseCase,
) *Handler {
	return &Handler{
		apiKey:            apiKey,
		createTaskUseCase: createTaskUseCase,
		getTaskUseCase:    getTaskUseCase,
	}
}
