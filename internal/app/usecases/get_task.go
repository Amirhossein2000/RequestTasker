package usecases

import (
	"context"

	"github.com/Amirhossein2000/RequestTasker/internal/app/services/logger"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/entities"

	"github.com/google/uuid"
)

type GetTaskUseCase struct {
	logger               logger.Logger
	taskRepository       entities.TaskRepository
	taskStatusRepository entities.TaskStatusRepository
	taskResultRepository entities.TaskResultRepository
}

func NewGetTaskUseCase(
	logger *logger.Logger,
	taskRepository entities.TaskRepository,
	taskStatusRepository entities.TaskStatusRepository,
	taskResultRepository entities.TaskResultRepository,
) GetTaskUseCase {
	return GetTaskUseCase{
		logger:               *logger,
		taskRepository:       taskRepository,
		taskStatusRepository: taskStatusRepository,
		taskResultRepository: taskResultRepository,
	}
}

func (u GetTaskUseCase) Execute(ctx context.Context, publicID uuid.UUID) (*entities.Task, *entities.TaskStatus, *entities.TaskResult, error) {
	task, err := u.taskRepository.GetByPublicID(ctx, publicID)
	if err != nil {
		return nil, nil, nil, err
	}

	taskStatus, err := u.taskStatusRepository.GetLatestByTaskID(ctx, task.ID())
	if err != nil {
		// TODO log
		return nil, nil, nil, err
	}

	if !taskStatus.HasResult() {
		return task, taskStatus, nil, nil
	}

	taskResult, err := u.taskResultRepository.GetByTaskID(ctx, task.ID())
	if err != nil {
		// TODO log
		return nil, nil, nil, err
	}

	return task, taskStatus, taskResult, nil
}
