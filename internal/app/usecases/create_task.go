package usecases

import (
	"RequestTasker/internal/app/services/logger"
	"RequestTasker/internal/domain/common"
	"RequestTasker/internal/domain/entities"
	"context"

	"github.com/google/uuid"
)

//go:generate mockery --name Tasker --structname TaskerMock --output ../../mocks/
type Tasker interface {
	RegisterTask(ctx context.Context, task entities.Task) error
}

type CreateTaskUseCase struct {
	logger               logger.Logger
	taskRepository       entities.TaskRepository
	tasker               Tasker
	taskStatusRepository entities.TaskStatusRepository
}

func NewCreateTaskUseCase(
	logger logger.Logger,
	taskRepository entities.TaskRepository,
	taskStatusRepository entities.TaskStatusRepository,
	requestTasker Tasker,
) CreateTaskUseCase {
	return CreateTaskUseCase{
		logger:               logger,
		taskRepository:       taskRepository,
		taskStatusRepository: taskStatusRepository,
		tasker:               requestTasker,
	}
}

func (u CreateTaskUseCase) Execute(ctx context.Context, task entities.Task) (uuid.UUID, error) {
	createdTask, err := u.taskRepository.Create(ctx, task)
	if err != nil {
		// TODO log
		return uuid.Nil, common.ErrInternal
	}

	status := entities.NewTaskStatus(createdTask.ID(), common.StatusNEW)
	_, err = u.taskStatusRepository.Create(ctx, status)
	if err != nil {
		// TODO log
		return uuid.Nil, common.ErrInternal
	}

	err = u.tasker.RegisterTask(ctx, *createdTask)
	if err != nil {
		// TODO log
		return uuid.Nil, common.ErrInternal
	}

	return createdTask.PublicID(), nil
}
