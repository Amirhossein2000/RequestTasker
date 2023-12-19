package entities

import (
	"RequestTasker/internal/domain/common"
	"context"
	"slices"
	"time"
)

//go:generate mockery --name TaskStatusRepository --structname TaskStatusRepositoryMock --output ../../mocks/
type TaskStatusRepository interface {
	Create(ctx context.Context, taskStatus TaskStatus) (*TaskStatus, error)
	GetLatestByTaskID(ctx context.Context, taskID int64) (*TaskStatus, error)
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

func (ts TaskStatus) HasResult() bool {
	return slices.Contains(
		[]string{
			common.StatusDONE,
			common.StatusERROR,
		},
		ts.status)
}
