package mysql

import (
	"RequestTasker/internal/domian/common"
	"RequestTasker/internal/domian/entities"
	"RequestTasker/internal/pkg/integration"
	"RequestTasker/internal/pkg/test"
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTaskStatusRepository(t *testing.T) {
	Convey("TaskStatusRepository INSERT and SELECT queries", t, func() {
		session, tearDown, err := integration.SetupMySQLContainer()
		So(err, ShouldBeNil)
		defer tearDown()

		repo := NewTaskStatusRepository(session, common.TaskStatusTable)

		task := test.NewTestTask()
		taskRepo := NewTaskRepository(session, common.TaskTable)
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
