package mysql

import (
	"RequestTasker/internal/domian/common"
	"RequestTasker/internal/domian/entities"
	"RequestTasker/internal/pkg/integration"
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTaskStatusRepository(t *testing.T) {
	Convey("TaskStatusRepository INSERT and SELECT queries", t, func() {
		conn, tearDown, err := integration.SetupMySQLContainer()
		So(err, ShouldBeNil)
		defer tearDown()

		session := conn.NewSession(nil)
		repo := NewTaskStatusRepository(session, "task_statuses")

		task := entities.NewTask(
			"https://example.com",
			"GET",
			map[string]interface{}{"Authorization": "Bearer token"},
			map[string]interface{}{"key": "value"},
		)
		taskRepo := NewTaskRepository(session, "tasks")
		createdTask, err := taskRepo.Create(context.Background(), task)
		So(err, ShouldBeNil)

		taskStatus := entities.NewTaskStatus(
			createdTask.ID(),
			common.StatusNEW,
		)

		Convey("Insert new task and check created task", func() {
			createdTaskStatus, err := repo.Create(context.Background(), taskStatus)
			So(err, ShouldBeNil)

			So(createdTaskStatus.ID(), ShouldNotEqual, int64(0))
			So(createdTaskStatus.TaskID(), ShouldEqual, taskStatus.TaskID())
			So(createdTaskStatus.Status(), ShouldEqual, taskStatus.Status())

			Convey("TaskStatusRepository.GetLatestByTaskID()", func() {
				foundTaskStatus, err := repo.GetLatestByTaskID(context.Background(), createdTaskStatus.TaskID())
				So(err, ShouldBeNil)

				So(createdTaskStatus.ID(), ShouldEqual, foundTaskStatus.ID())
				So(createdTaskStatus.TaskID(), ShouldEqual, foundTaskStatus.TaskID())
				So(createdTaskStatus.Status(), ShouldEqual, foundTaskStatus.Status())
			})
		})
	})
}