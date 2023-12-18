package mysql

import (
	"RequestTasker/internal/domian/entities"
	"RequestTasker/internal/pkg/integration"
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTaskResultRepository(t *testing.T) {
	Convey("TaskResultRepository INSERT and SELECT queries", t, func() {
		conn, tearDown, err := integration.SetupMySQLContainer()
		So(err, ShouldBeNil)
		defer tearDown()

		session := conn.NewSession(nil)
		repo := NewTaskResultRepository(session, "task_results")

		task := entities.NewTask(
			"https://example.com",
			"GET",
			map[string]interface{}{"Authorization": "Bearer token"},
			map[string]interface{}{"key": "value"},
		)
		taskRepo := NewTaskRepository(session, "tasks")
		createdTask, err := taskRepo.Create(context.Background(), task)
		So(err, ShouldBeNil)

		TaskResult := entities.NewTaskResult(
			createdTask.ID(),
			200,
			map[string]interface{}{
				"OK": "OK",
			},
			10,
		)

		Convey("Insert new task and check created task", func() {
			createdTaskResult, err := repo.Create(context.Background(), TaskResult)
			So(err, ShouldBeNil)

			So(createdTaskResult.ID(), ShouldNotEqual, int64(0))
			So(createdTaskResult.TaskID(), ShouldEqual, TaskResult.TaskID())
			So(createdTaskResult.StatusCode(), ShouldEqual, TaskResult.StatusCode())
			So(createdTaskResult.Headers(), ShouldEqual, TaskResult.Headers())
			So(createdTaskResult.Length(), ShouldEqual, TaskResult.Length())

			Convey("TaskResultRepository.GetByPublicID()", func() {
				foundTaskResult, err := repo.GetByTaskID(context.Background(), createdTaskResult.TaskID())
				So(err, ShouldBeNil)

				So(createdTaskResult.ID(), ShouldEqual, foundTaskResult.ID())
				So(createdTaskResult.TaskID(), ShouldEqual, foundTaskResult.TaskID())
				So(createdTaskResult.StatusCode(), ShouldEqual, TaskResult.StatusCode())
				So(createdTaskResult.Headers(), ShouldEqual, TaskResult.Headers())
				So(createdTaskResult.Length(), ShouldEqual, TaskResult.Length())
			})
		})
	})
}