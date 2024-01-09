package usecases

import (
	"context"

	"github.com/Amirhossein2000/RequestTasker/internal/domain/entities"

	"github.com/google/uuid"
)

type GetTaskUseCase struct {
	taskRepository       entities.TaskRepository
	taskStatusRepository entities.TaskStatusRepository
	taskResultRepository entities.TaskResultRepository
}

func NewGetTaskUseCase(
	taskRepository entities.TaskRepository,
	taskStatusRepository entities.TaskStatusRepository,
	taskResultRepository entities.TaskResultRepository,
) GetTaskUseCase {
	return GetTaskUseCase{
		taskRepository:       taskRepository,
		taskStatusRepository: taskStatusRepository,
		taskResultRepository: taskResultRepository,
	}
}

// GetTaskUseCase returns the latest status and results of a tasks by task.PublicID
func (u GetTaskUseCase) Execute(ctx context.Context, publicID uuid.UUID) (*entities.Task, *entities.TaskStatus, *entities.TaskResult, error) {
	task, err := u.taskRepository.GetByPublicID(ctx, publicID)
	if err != nil {
		return nil, nil, nil, err
	}

	taskStatus, err := u.taskStatusRepository.GetLatestByTaskID(ctx, task.ID())
	if err != nil {
		return nil, nil, nil, err
	}

	if !taskStatus.HasResult() {
		return task, taskStatus, nil, nil
	}

	taskResult, err := u.taskResultRepository.GetByTaskID(ctx, task.ID())
	if err != nil {
		return nil, nil, nil, err
	}

	return task, taskStatus, taskResult, nil
}
