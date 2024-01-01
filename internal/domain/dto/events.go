package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type TaskEvent struct {
	PublicID uuid.UUID `json:"public_id"`
}

func NewTaskEvent(data []byte) (*TaskEvent, error) {
	taskEvent := &TaskEvent{}
	err := json.Unmarshal(data, taskEvent)
	return taskEvent, err
}

func (t *TaskEvent) Serialize() ([]byte, error) {
	return json.Marshal(t)
}
