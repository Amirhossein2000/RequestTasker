package mysql

import (
	"RequestTasker/internal/domian/common"
	"RequestTasker/internal/domian/entities"
	"context"
	"fmt"
	"time"

	"github.com/gocraft/dbr/v2"
)

type TaskStatusRow struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	TaskID    int64     `db:"task_id"`
	Status    string    `db:"status"`
}

func (r *TaskStatusRow) ConvertToEntity() (*entities.TaskStatus, error) {
	taskStatus := entities.BuildTaskStatus(
		r.ID,
		r.CreatedAt,
		r.TaskID,
		r.Status,
	)

	return &taskStatus, nil
}

type TaskStatusRepository struct {
	session   *dbr.Session
	tableName string
}

func NewTaskStatusRepository(session *dbr.Session, tableName string) *TaskStatusRepository {
	return &TaskStatusRepository{
		session:   session,
		tableName: tableName,
	}
}

func (r *TaskStatusRepository) Create(ctx context.Context, taskStatus entities.TaskStatus) (*entities.TaskStatus, error) {
	result, err := r.session.InsertInto(r.tableName).
		Columns(
			"task_id",
			"status",
		).
		Values(
			taskStatus.TaskID(),
			taskStatus.Status(),
		).
		ExecContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create taskStatus: %w", err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	createdTaskStatus := entities.BuildTaskStatus(
		lastInsertID,
		taskStatus.CreatedAt(),
		taskStatus.TaskID(),
		taskStatus.Status(),
	)

	return &createdTaskStatus, nil
}

func (r *TaskStatusRepository) GetLatestByTaskID(ctx context.Context, taskID int64) (*entities.TaskStatus, error) {
	taskStatusRow := &TaskStatusRow{}

	err := r.session.Select("*").
		From(r.tableName).
		Where(dbr.Eq("task_id", taskID)).
		Limit(1).
		LoadOneContext(ctx, taskStatusRow)

	switch {
	case err == dbr.ErrNotFound:
		return nil, common.NotFoundError
	case err != nil:
		return nil, fmt.Errorf("failed to get task by public ID: %w", err)
	}

	return taskStatusRow.ConvertToEntity()
}
