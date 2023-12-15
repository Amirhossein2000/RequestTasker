package entities

import (
	"time"
)

//go:generate mockery --name TaskResultRepository --structname TaskResultRepositoryMock --output ../../mocks/ 
type TaskResultRepository interface {
	Create(TaskResult) error
	GetByTaskID(taskID int64) (TaskResult, error)
}

type TaskResult struct {
	id         int64
	createdAt  time.Time
	taskID     int64
	statusCode int
	headers    map[string]interface{}
	length     int64
}

func NewTaskResult(
	taskID int64,
	statusCode int,
	headers map[string]interface{},
	length int64,
) TaskResult {
	return TaskResult{
		createdAt:  time.Now(),
		taskID:     taskID,
		statusCode: statusCode,
		headers:    headers,
		length:     length,
	}
}

func BuildTaskResult(
	id int64,
	createdAt time.Time,
	taskID int64,
	statusCode int,
	headers map[string]interface{},
	length int64,
) TaskResult {
	return TaskResult{
		id:         id,
		createdAt:  createdAt,
		taskID:     taskID,
		statusCode: statusCode,
		headers:    headers,
		length:     length,
	}
}

func (tr TaskResult) ID() int64 {
	return tr.id
}

func (tr TaskResult) CreatedAt() time.Time {
	return tr.createdAt
}

func (tr TaskResult) TaskID() int64 {
	return tr.taskID
}

func (tr TaskResult) StatusCode() int {
	return tr.statusCode
}

func (tr TaskResult) Headers() map[string]interface{} {
	return tr.headers
}

func (tr TaskResult) Length() int64 {
	return tr.length
}
