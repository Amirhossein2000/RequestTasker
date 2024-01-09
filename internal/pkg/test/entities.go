package test

import "github.com/Amirhossein2000/RequestTasker/internal/domain/entities"

func NewTestTask() entities.Task {
	return entities.NewTask(
		"https://example.com",
		"GET",
		map[string]string{
			"Authorization": "Bearer token test",
		},
		`
			{
				"test":	"test"
			}
		`,
	)
}
