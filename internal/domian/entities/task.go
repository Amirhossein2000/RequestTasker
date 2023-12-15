package entities

import (
	"time"

	"github.com/google/uuid"
)

type TaskRepository interface {
	Create(Task) error
	GetByPublicID(publicId uuid.UUID) (Task, error)
}

type Task struct {
	id        int64
	createdAt time.Time
	publicID  uuid.UUID

	url     string
	method  string
	headers map[string]interface{}
	body    map[string]interface{}
}

func NewTask(
	url string,
	method string,
	headers map[string]interface{},
	body map[string]interface{},
) Task {
	return Task{
		createdAt: time.Now(),
		publicID:  uuid.New(),
		url:       url,
		method:    method,
		headers:   headers,
		body:      body,
	}
}

func BuildTask(
	id int64,
	createdAt time.Time,
	publicID uuid.UUID,
	url string,
	method string,
	headers map[string]interface{},
	body map[string]interface{},
) Task {
	return Task{
		id:        id,
		createdAt: createdAt,
		publicID:  publicID,
		url:       url,
		method:    method,
		headers:   headers,
		body:      body,
	}
}

func (t Task) ID() int64 {
	return t.id
}

func (t Task) CreatedAt() time.Time {
	return t.createdAt
}

func (t Task) Url() string {
	return t.url
}

func (t Task) Method() string {
	return t.method
}

func (t Task) Headers() map[string]interface{} {
	return t.headers
}

func (t Task) Body() map[string]interface{} {
	return t.body
}
