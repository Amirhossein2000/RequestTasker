package usecases

import (
	"RequestTasker/internal/app/services/logger"
	"RequestTasker/internal/domian/common"
	"RequestTasker/internal/domian/entities"

	"github.com/google/uuid"
)

type GetTaskUseCase struct {
	logger               logger.Logger
	taskRepository       entities.TaskRepository
	taskStatusRepository entities.TaskStatusRepository
	taskResultRepository entities.TaskResultRepository
}

func NewGetTaskUseCase(
	logger logger.Logger,
	taskRepository entities.TaskRepository,
	taskStatusRepository entities.TaskStatusRepository,
	taskResultRepository entities.TaskResultRepository,
) GetTaskUseCase {
	return GetTaskUseCase{
		logger:               logger,
		taskRepository:       taskRepository,
		taskStatusRepository: taskStatusRepository,
		taskResultRepository: taskResultRepository,
	}
}

func (u GetTaskUseCase) Execute(publicID uuid.UUID) (*entities.Task, *entities.TaskStatus, *entities.TaskResult, error) {
	task, err := u.taskRepository.GetByPublicID(publicID)
	if err != nil {
		// TODO log
		return nil, nil, nil, common.InternalError
	}

	taskStatus, err := u.taskStatusRepository.GetLatestByTaskID(task.ID())
	if err != nil {
		// TODO log
		return nil, nil, nil, common.InternalError
	}

	if !taskStatus.HasResult() {
		return &task, &taskStatus, nil, nil
	}

	taskResult, err := u.taskResultRepository.GetByTaskID(task.ID())
	if err != nil {
		// TODO log
		return nil, nil, nil, common.InternalError
	}

	return &task, &taskStatus, &taskResult, nil
}
