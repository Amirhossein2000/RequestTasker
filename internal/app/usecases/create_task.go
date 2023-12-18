package usecases

import (
	"RequestTasker/internal/app/services/logger"
	"RequestTasker/internal/domian/common"
	"RequestTasker/internal/domian/entities"
	"context"

	"github.com/google/uuid"
)

//go:generate mockery --name RequestTasker --structname RequestTaskerMock --output ../../mocks/
type RequestTasker interface {
	RegisterTask(ctx context.Context, task entities.Task) error
}

type CreateTaskUseCase struct {
	logger               logger.Logger
	taskRepository       entities.TaskRepository
	requestTasker        RequestTasker
	taskStatusRepository entities.TaskStatusRepository
}

func NewCreateTaskUseCase(
	logger logger.Logger,
	taskRepository entities.TaskRepository,
	taskStatusRepository entities.TaskStatusRepository,
	requestTasker RequestTasker,
) CreateTaskUseCase {
	return CreateTaskUseCase{
		logger:               logger,
		taskRepository:       taskRepository,
		taskStatusRepository: taskStatusRepository,
		requestTasker:        requestTasker,
	}
}

func (u CreateTaskUseCase) Execute(ctx context.Context, task entities.Task) (uuid.UUID, error) {
	createdTask, err := u.taskRepository.Create(ctx, task)
	if err != nil {
		// TODO log
		return uuid.Nil, common.InternalError
	}

	status := entities.NewTaskStatus(createdTask.ID(), common.StatusNEW)
	_, err = u.taskStatusRepository.Create(ctx, status)
	if err != nil {
		// TODO log
		return uuid.Nil, common.InternalError
	}

	err = u.requestTasker.RegisterTask(ctx, *createdTask)
	if err != nil {
		// TODO log
		return uuid.Nil, common.InternalError
	}

	return createdTask.PublicID(), nil
}
