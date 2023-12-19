package mysql

import (
	"RequestTasker/internal/domian/common"
	"RequestTasker/internal/domian/entities"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocraft/dbr/v2"
)

type TaskResultRow struct {
	ID         int64     `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	TaskID     int64     `db:"task_id"`
	StatusCode int       `db:"status_code"`
	Headers    string    `db:"headers"`
	Length     int64     `db:"length"`
}

func (r *TaskResultRow) ConvertToEntity() (*entities.TaskResult, error) {
	headers := make(map[string]string)
	err := json.Unmarshal([]byte(r.Headers), &headers)
	if err != nil {
		return nil, err
	}

	TaskResult := entities.BuildTaskResult(
		r.ID,
		r.CreatedAt,
		r.TaskID,
		r.StatusCode,
		headers,
		r.Length,
	)

	return &TaskResult, nil
}

type TaskResultRepository struct {
	session   *dbr.Session
	tableName string
}

func NewTaskResultRepository(session *dbr.Session, tableName string) *TaskResultRepository {
	return &TaskResultRepository{
		session:   session,
		tableName: tableName,
	}
}

func (r *TaskResultRepository) Create(ctx context.Context, TaskResult entities.TaskResult) (*entities.TaskResult, error) {
	headers, err := json.Marshal(TaskResult.Headers())
	if err != nil {
		return nil, err
	}

	result, err := r.session.InsertInto(r.tableName).
		Columns(
			"task_id",
			"status_code",
			"headers",
			"length",
		).
		Values(
			TaskResult.TaskID(),
			TaskResult.StatusCode(),
			string(headers),
			TaskResult.Length(),
		).
		ExecContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create TaskResult: %w", err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last TaskResult insert ID: %w", err)
	}

	createdTaskResult := entities.BuildTaskResult(
		lastInsertID,
		TaskResult.CreatedAt(),
		TaskResult.TaskID(),
		TaskResult.StatusCode(),
		TaskResult.Headers(),
		TaskResult.Length(),
	)

	return &createdTaskResult, nil
}

func (r *TaskResultRepository) GetByTaskID(ctx context.Context, taskID int64) (*entities.TaskResult, error) {
	TaskResultRow := &TaskResultRow{}

	err := r.session.Select("*").
		From(r.tableName).
		Where(dbr.Eq("task_id", taskID)).
		Limit(1).
		LoadOneContext(ctx, TaskResultRow)

	switch {
	case err == dbr.ErrNotFound:
		return nil, common.ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("failed to GetByTaskID: %w", err)
	}

	return TaskResultRow.ConvertToEntity()
}
