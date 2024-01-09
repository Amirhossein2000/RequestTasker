package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Amirhossein2000/RequestTasker/internal/domain/common"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/entities"
	"github.com/samber/lo"

	"github.com/gocraft/dbr/v2"
	"github.com/google/uuid"
)

type TaskRow struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	PublicID  string    `db:"public_id"`

	Url     string  `db:"url"`
	Method  string  `db:"method"`
	Headers *string `db:"headers"`
	Body    *string `db:"body"`
}

func (r *TaskRow) ConvertToEntity() (*entities.Task, error) {
	headers := make(map[string]string)

	if r.Headers != nil {
		err := json.Unmarshal([]byte(*r.Headers), &headers)
		if err != nil {
			return nil, err
		}
	} else {
		headers = nil
	}

	publicID, err := uuid.Parse(r.PublicID)
	if err != nil {
		return nil, err
	}

	body := ""
	if r.Body != nil {
		body = *r.Body
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

type TaskRepository struct {
	session   *dbr.Session
	tableName string
}

func NewTaskRepository(conn *dbr.Connection, tableName string) *TaskRepository {
	return &TaskRepository{
		session:   conn.NewSession(nil),
		tableName: tableName,
	}
}

func (r *TaskRepository) Create(ctx context.Context, task entities.Task) (*entities.Task, error) {
	var headers *string

	if task.Headers() != nil {
		data, err := json.Marshal(task.Headers())
		if err != nil {
			return nil, err
		}
		headers = lo.ToPtr(string(data))
	} else {
		headers = nil
	}

	var body *string
	if taskBody := task.Body(); taskBody != "" {
		body = &taskBody
	} else {
		body = nil
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
			headers,
			body,
		).
		ExecContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last task insert ID: %w", err)
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

func (r *TaskRepository) GetByPublicID(ctx context.Context, publicID uuid.UUID) (*entities.Task, error) {
	taskRow := &TaskRow{}

	err := r.session.Select("*").
		From(r.tableName).
		Where(dbr.Eq("public_id", publicID)).
		Limit(1).
		LoadOneContext(ctx, taskRow)

	switch {
	case err == dbr.ErrNotFound:
		return nil, common.ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("failed to GetByPublicID: %w", err)
	}

	return taskRow.ConvertToEntity()
}
