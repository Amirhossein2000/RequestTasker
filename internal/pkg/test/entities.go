package test

import "RequestTasker/internal/domian/entities"

func NewTestTask() entities.Task {
	return entities.NewTask(
		"https://example.com",
		"GET",
		map[string]string{"Authorization": "Bearer token test"},
		`
			{
				"test":	"test"
			}
		`,
	)
}
