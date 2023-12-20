package handlers

import "RequestTasker/internal/app/usecases"

type Handler struct {
	createTaskUseCase usecases.CreateTaskUseCase
	getTaskUseCase    usecases.GetTaskUseCase
}

func NewHandler(
	createTaskUseCase usecases.CreateTaskUseCase,
	getTaskUseCase usecases.GetTaskUseCase,
) *Handler {
	return &Handler{
		createTaskUseCase: createTaskUseCase,
		getTaskUseCase:    getTaskUseCase,
	}
}
