package tasker

import (
	"RequestTasker/internal/domain/common"
	"RequestTasker/internal/domain/dto"
	"RequestTasker/internal/domain/entities"
	"bytes"
	"context"
	"net/http"
	"strings"
)

//go:generate mockery --name TaskEventRepository --structname TaskEventRepositoryMock --output ../../../mocks/
type TaskEventRepository interface {
	Read(ctx context.Context) ([]byte, error)
	Write(ctx context.Context, value []byte) error
}

type Tasker struct {
	taskEventRepository  TaskEventRepository
	taskRepository       entities.TaskRepository
	taskStatusRepository entities.TaskStatusRepository
	taskResultRepository entities.TaskResultRepository
	httpClient           *http.Client
	in                   chan int64
	out                  chan int64
}

func NewTasker(
	taskEventRepository TaskEventRepository,
	taskRepository entities.TaskRepository,
	taskStatusRepository entities.TaskStatusRepository,
	taskResultRepository entities.TaskResultRepository,
) *Tasker {
	return &Tasker{
		taskEventRepository:  taskEventRepository,
		taskRepository:       taskRepository,
		taskStatusRepository: taskStatusRepository,
		taskResultRepository: taskResultRepository,
		in:                   make(chan int64, 128),
		out:                  make(chan int64, 128),
	}
}

func (t *Tasker) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case taskID := <-t.in:
			event := dto.TaskEvent{
				ID: taskID,
			}
			data, err := event.Serialize()
			if err != nil {
				return err
			}
			err = t.taskEventRepository.Write(ctx, data)
			if err != nil {
				return err
			}

			taskStatus := entities.NewTaskStatus(
				taskID,
				common.StatusIN_PROGRESS,
			)

			_, err = t.taskStatusRepository.Create(ctx, taskStatus)
			if err != nil {
				return err
			}

		default:
			data, err := t.taskEventRepository.Read(ctx)
			if err != nil {
				return err
			}

			event, err := dto.NewTaskEvent(data)
			if err != nil {
				return err
			}

			err = t.sendTask(ctx, event.ID)
			if err != nil {
				return err
			}
		}
	}
}

func (t *Tasker) RegisterTask(ctx context.Context, task entities.Task) error {
	select {
	case t.in <- task.ID():
		return nil

	default:
		return common.ErrChannelDeadlock
	}
}

func (t *Tasker) sendTask(ctx context.Context, taskID int64) error {
	task, err := t.taskRepository.Get(ctx, taskID)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer([]byte(task.Body()))

	req, err := http.NewRequestWithContext(ctx, task.Method(), task.Url(), body)
	if err != nil {
		return err
	}

	for key, value := range task.Headers() {
		req.Header.Set(key, value)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) < 1 {
			headers[k] = "EMPTY"
		} else if len(v) == 1 {
			headers[k] = v[0]
		} else {
			headers[k] = strings.Join(v, "|")
		}
	}

	result := entities.NewTaskResult(
		task.ID(),
		resp.StatusCode,
		headers,
		resp.ContentLength,
	)

	_, err = t.taskResultRepository.Create(ctx, result)
	if err != nil {
		return err
	}

	taskStatus := entities.NewTaskStatus(
		task.ID(),
		common.StatusDONE,
	)

	_, err = t.taskStatusRepository.Create(ctx, taskStatus)
	if err != nil {
		return err
	}

	return nil
}

// TODO, send all of the NEW tasks?
