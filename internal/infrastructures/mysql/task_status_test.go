package mysql

import (
	"context"
	"testing"

	"github.com/Amirhossein2000/RequestTasker/internal/domain/common"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/entities"
	"github.com/Amirhossein2000/RequestTasker/internal/pkg/integration"
	"github.com/Amirhossein2000/RequestTasker/internal/pkg/test"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTaskStatusRepository(t *testing.T) {
	conn, tearDown, err := integration.SetupMySQLContainer()
	if err != nil {
		t.Fatal(err)
	}
	defer tearDown()

	repo := NewTaskStatusRepository(conn, common.TaskStatusTable)

	task := test.NewTestTask()
	taskRepo := NewTaskRepository(conn, common.TaskTable)
	createdTask, err := taskRepo.Create(context.Background(), task)
	if err != nil {
		t.Fatal(err)
	}

	taskStatusNEW := entities.NewTaskStatus(
		createdTask.ID(),
		common.StatusNEW,
	)

	Convey("TaskStatusRepository INSERT and SELECT queries", t, func() {
		Convey("Insert new task and check created task", func() {
			createdTaskStatus, err := repo.Create(context.Background(), taskStatusNEW)
			So(err, ShouldBeNil)

			So(createdTaskStatus.ID(), ShouldNotEqual, int64(0))
			So(createdTaskStatus.TaskID(), ShouldEqual, taskStatusNEW.TaskID())
			So(createdTaskStatus.Status(), ShouldEqual, taskStatusNEW.Status())

			taskStatusDONE := entities.NewTaskStatus(
				createdTask.ID(),
				common.StatusDONE,
			)
			createdTaskStatus, err = repo.Create(context.Background(), taskStatusDONE)
			So(err, ShouldBeNil)
			So(createdTaskStatus.ID(), ShouldNotEqual, int64(0))
			So(createdTaskStatus.TaskID(), ShouldEqual, taskStatusDONE.TaskID())
			So(createdTaskStatus.Status(), ShouldEqual, taskStatusDONE.Status())

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
