package usecases

import (
	"context"

	"github.com/Amirhossein2000/RequestTasker/internal/domain/common"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/entities"

	"github.com/google/uuid"
)

//go:generate mockery --name Tasker --structname TaskerMock --output ../../mocks/
type Tasker interface {
	Process(ctx context.Context, task entities.Task) error
}

type CreateTaskUseCase struct {
	taskRepository       entities.TaskRepository
	tasker               Tasker
	taskStatusRepository entities.TaskStatusRepository
}

func NewCreateTaskUseCase(
	taskRepository entities.TaskRepository,
	taskStatusRepository entities.TaskStatusRepository,
	tasker Tasker,
) CreateTaskUseCase {
	return CreateTaskUseCase{
		taskRepository:       taskRepository,
		taskStatusRepository: taskStatusRepository,
		tasker:               tasker,
	}
}

// CreateTaskUseCase creates entities and process the new task
func (u CreateTaskUseCase) Execute(ctx context.Context, task entities.Task) (uuid.UUID, error) {
	createdTask, err := u.taskRepository.Create(ctx, task)
	if err != nil {
		return uuid.Nil, err
	}

	status := entities.NewTaskStatus(createdTask.ID(), common.StatusNew)
	_, err = u.taskStatusRepository.Create(ctx, status)
	if err != nil {
		return uuid.Nil, err
	}

	err = u.tasker.Process(ctx, *createdTask)
	if err != nil {
		return uuid.Nil, err
	}

	return createdTask.PublicID(), nil
}
