package dto

import (
	"encoding/json"
)

type TaskEvent struct {
	ID int64
}

func NewTaskEvent(data []byte) (*TaskEvent, error) {
	taskEvent := &TaskEvent{}
	err := json.Unmarshal(data, taskEvent)
	return taskEvent, err
}

func (t *TaskEvent) Serialize() ([]byte, error) {
	return json.Marshal(t)
}
