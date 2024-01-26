/*
	Main jobs in this package:
		1. Emit tasks to the message broker for distribution.
		2. Acquire and process tasks from the message broker.
		3. Dispatch tasks to third-party endpoints.
		4. Archive task results and their statuses in the database.
*/

package tasker

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Amirhossein2000/RequestTasker/internal/app/services/logger"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/common"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/dto"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/entities"
	"github.com/google/uuid"
)

//go:generate mockery --name TaskEventRepository --structname TaskEventRepositoryMock --output ../../../mocks/
type TaskEventRepository interface {
	Read(ctx context.Context) ([]byte, error)
	Write(ctx context.Context, value []byte) error
}

type Tasker struct {
	logger               *logger.Logger
	taskEventRepository  TaskEventRepository
	taskRepository       entities.TaskRepository
	taskStatusRepository entities.TaskStatusRepository
	taskResultRepository entities.TaskResultRepository
	httpClient           *http.Client
	produceChan          chan *entities.Task
	consumeChan          chan []byte
}

func NewTasker(
	logger *logger.Logger,
	taskEventRepository TaskEventRepository,
	taskRepository entities.TaskRepository,
	taskStatusRepository entities.TaskStatusRepository,
	taskResultRepository entities.TaskResultRepository,
	httpClient *http.Client,
) *Tasker {
	return &Tasker{
		logger:               logger,
		taskEventRepository:  taskEventRepository,
		taskRepository:       taskRepository,
		taskStatusRepository: taskStatusRepository,
		taskResultRepository: taskResultRepository,
		httpClient:           httpClient,
		produceChan:          make(chan *entities.Task, 128),
		consumeChan:          make(chan []byte, 128),
	}
}

func (t *Tasker) Start(ctx context.Context) {
	for i := 0; i < 10; i++ {
		go t.eventHandler(ctx)
	}

	go t.produce(ctx)
	go t.consume(ctx)
}

func (t *Tasker) Shutdown() {
	close(t.produceChan)
	close(t.consumeChan)
}

func (t *Tasker) produce(ctx context.Context) {
	for task := range t.produceChan {
		event := dto.TaskEvent{
			PublicID: task.PublicID(),
		}
		data, err := event.Serialize()
		if err != nil {
			t.logger.Error("event Serialize failed",
				"error", err,
			)
			continue
		}
		err = t.taskEventRepository.Write(ctx, data)
		if err != nil {
			t.logger.Error("write event failed",
				"error", err,
			)
			continue
		}

		taskStatus := entities.NewTaskStatus(
			task.ID(),
			common.StatusInProcess,
		)

		_, err = t.taskStatusRepository.Create(ctx, taskStatus)
		if err != nil {
			t.logger.Error("create taskStatus failed",
				"error", err,
			)
			continue
		}
	}
}

func (t *Tasker) consume(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, err := t.taskEventRepository.Read(ctx)
			if err != nil {
				t.logger.Error("read event failed",
					"error", err,
				)
				continue
			}

			t.consumeChan <- data
		}
	}
}

func (t *Tasker) Process(ctx context.Context, task entities.Task) error {
	select {
	case t.produceChan <- &task:
		return nil

	default:
		return common.ErrChannelDeadlock
	}
}

func (t *Tasker) eventHandler(ctx context.Context) {
	for data := range t.consumeChan {
		event, err := dto.NewTaskEvent(data)
		if err != nil {
			t.logger.Error("decentralizing event failed",
				"error", err,
			)
			continue
		}

		err = t.sendTask(ctx, event.PublicID)
		if err != nil {
			t.logger.Error("sending task failed",
				"error", err,
			)
			continue
		}
	}
}

func (t *Tasker) sendTask(ctx context.Context, taskPublicID uuid.UUID) error {
	task, err := t.taskRepository.GetByPublicID(ctx, taskPublicID)
	if err != nil {
		return err
	}

	var req *http.Request
	if task.Body() != "" {
		body := bytes.NewBuffer([]byte(task.Body()))
		req, err = http.NewRequestWithContext(ctx, task.Method(), task.Url(), body)
		if err != nil {
			return err
		}
	} else {
		req, err = http.NewRequestWithContext(ctx, task.Method(), task.Url(), nil)
		if err != nil {
			return err
		}
	}

	for key, value := range task.Headers() {
		req.Header.Set(key, value)
	}

	resp, err := t.httpClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		taskStatus := entities.NewTaskStatus(
			task.ID(),
			common.StatusError,
		)

		_, dbErr := t.taskStatusRepository.Create(ctx, taskStatus)
		if err != nil {
			return fmt.Errorf("send request error: %w, db error: %w", err, dbErr)
		}
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

	var taskStatus entities.TaskStatus
	if resp.StatusCode > 199 && resp.StatusCode < 300 {
		taskStatus = entities.NewTaskStatus(
			task.ID(),
			common.StatusDone,
		)
	} else {
		taskStatus = entities.NewTaskStatus(
			task.ID(),
			common.StatusError,
		)
	}

	_, err = t.taskStatusRepository.Create(ctx, taskStatus)
	if err != nil {
		return err
	}

	return nil
}
