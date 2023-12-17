package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"RequestTasker/internal/domian/common"
	"RequestTasker/internal/domian/entities"

	"github.com/gocraft/dbr/v2"
	"github.com/google/uuid"
)

type TaskRow struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	PublicID  string    `db:"public_id"`

	Url     string `db:"url"`
	Method  string `db:"method"`
	Headers string `db:"headers"`
	Body    string `db:"body"`
}

func (r *TaskRow) ConvertToEntity() (*entities.Task, error) {
	headers := make(map[string]interface{})
	body := make(map[string]interface{})

	err := json.Unmarshal([]byte(r.Headers), &headers)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(r.Body), &body)
	if err != nil {
		return nil, err
	}

	publicID, err := uuid.Parse(r.PublicID)
	if err != nil {
		return nil, err
	}

	task := entities.BuildTask(
		r.ID,
		r.CreatedAt,
		publicID,
		r.Url,
		r.Method,
		headers,
		body,
	)

	return &task, nil
}

type MySQLTaskRepository struct {
	session   *dbr.Session
	tableName string
}

func NewMySQLTaskRepository(session *dbr.Session, tableName string) *MySQLTaskRepository {
	return &MySQLTaskRepository{
		session:   session,
		tableName: tableName,
	}
}

func (r *MySQLTaskRepository) Create(ctx context.Context, task entities.Task) (*entities.Task, error) {
	headers, err := json.Marshal(task.Headers())
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(task.Body())
	if err != nil {
		return nil, err
	}

	result, err := r.session.InsertInto(r.tableName).
		Columns(
			"public_id",
			"url",
			"method",
			"headers",
			"body",
		).
		Values(
			task.PublicID().String(),
			task.Url(),
			task.Method(),
			string(headers),
			string(body),
		).
		ExecContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	createdTask := entities.BuildTask(
		lastInsertID,
		task.CreatedAt(),
		task.PublicID(),
		task.Url(),
		task.Method(),
		task.Headers(),
		task.Body(),
	)

	return &createdTask, nil
}

func (r *MySQLTaskRepository) GetByPublicID(ctx context.Context, publicID uuid.UUID) (*entities.Task, error) {
	taskRow := &TaskRow{}

	err := r.session.Select("*").
		From(r.tableName).
		Where(dbr.Eq("public_id", publicID)).
		Limit(1).
		LoadOneContext(ctx, taskRow)

	switch {
	case err == dbr.ErrNotFound:
		return nil, common.NotFoundError
	case err != nil:
		return nil, fmt.Errorf("failed to get task by public ID: %w", err)
	}

	return taskRow.ConvertToEntity()
}
