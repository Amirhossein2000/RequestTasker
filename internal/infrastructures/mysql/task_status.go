package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/Amirhossein2000/RequestTasker/internal/domain/common"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/entities"

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

func NewTaskStatusRepository(conn *dbr.Connection, tableName string) *TaskStatusRepository {
	return &TaskStatusRepository{
		session:   conn.NewSession(nil),
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
		return nil, fmt.Errorf("failed to get last taskStatus insert ID: %w", err)
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
		OrderBy("id DESC").
		LoadOneContext(ctx, taskStatusRow)

	switch {
	case err == dbr.ErrNotFound:
		return nil, common.ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("failed to GetLatestByTaskID: %w", err)
	}

	return taskStatusRow.ConvertToEntity()
}
