package entities

import (
	"context"
	"time"

	"github.com/google/uuid"
)

//go:generate mockery --name TaskRepository --structname TaskRepositoryMock --output ../../mocks/
type TaskRepository interface {
	Create(ctx context.Context, task Task) (*Task, error)
	GetByPublicID(ctx context.Context, publicId uuid.UUID) (*Task, error)
}

type Task struct {
	id        int64
	createdAt time.Time
	publicID  uuid.UUID

	url     string
	method  string
	headers map[string]string
	body    string // TODO: make it pointer
}

func NewTask(
	url string,
	method string,
	headers map[string]string,
	body string,
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
	headers map[string]string,
	body string,
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

func (t Task) PublicID() uuid.UUID {
	return t.publicID
}

func (t Task) Url() string {
	return t.url
}

func (t Task) Method() string {
	return t.method
}

func (t Task) Headers() map[string]string {
	return t.headers
}

func (t Task) Body() string {
	return t.body
}
