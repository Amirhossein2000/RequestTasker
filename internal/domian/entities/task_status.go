package entities

import (
	"time"
)

//go:generate mockery --name TaskStatusRepository --structname TaskStatusRepositoryMock --output ../../mocks/ 
type TaskStatusRepository interface {
	Create(TaskStatus) error
	GetByTaskID(taskID int64) (TaskStatus, error)
}

type TaskStatus struct {
	id        int64
	createdAt time.Time
	taskID    int64
	status    string
}

func NewTaskStatus(
	taskID int64,
	status string,
) TaskStatus {
	return TaskStatus{
		taskID: taskID,
		status: status,
	}
}

func BuildTaskStatus(
	id int64,
	createdAt time.Time,
	taskID int64,
	status string,
) TaskStatus {
	return TaskStatus{
		id:        id,
		createdAt: createdAt,
		taskID:    taskID,
		status:    status,
	}
}

func (ts TaskStatus) ID() int64 {
	return ts.id
}

func (ts TaskStatus) CreatedAt() time.Time {
	return ts.createdAt
}

func (ts TaskStatus) Status() string {
	return ts.status
}

func (ts TaskStatus) TaskID() int64 {
	return ts.taskID
}
