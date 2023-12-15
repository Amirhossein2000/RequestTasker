package usecases

import (
	"RequestTasker/internal/app/services/logger"
	"RequestTasker/internal/domian/common"
	"RequestTasker/internal/domian/entities"

	"github.com/google/uuid"
)

//go:generate mockery --name RequestTasker --structname RequestTaskerMock --output ../../mocks/
type RequestTasker interface {
	RegisterTask(entities.Task) error
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

func (u CreateTaskUseCase) Execute(task entities.Task) (uuid.UUID, error) {
	task, err := u.taskRepository.Create(task)
	if err != nil {
		// TODO log
		return uuid.Nil, common.InternalError
	}

	status := entities.NewTaskStatus(task.ID(), common.StatusNEW)
	err = u.taskStatusRepository.Create(status)
	if err != nil {
		// TODO log
		return uuid.Nil, common.InternalError
	}

	err = u.requestTasker.RegisterTask(task)
	if err != nil {
		// TODO log
		return uuid.Nil, common.InternalError
	}

	return task.PublicID(), nil
}
