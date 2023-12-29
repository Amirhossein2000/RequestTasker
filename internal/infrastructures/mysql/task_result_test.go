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

func TestTaskResultRepository(t *testing.T) {
	conn, tearDown, err := integration.SetupMySQLContainer()
	if err != nil {
		t.Fatal(err)
	}
	defer tearDown()

	repo := NewTaskResultRepository(conn, common.TaskResultTable)

	task := test.NewTestTask()
	taskRepo := NewTaskRepository(conn, common.TaskTable)
	createdTask, err := taskRepo.Create(context.Background(), task)
	if err != nil {
		t.Fatal(err)
	}

	TaskResult := entities.NewTaskResult(
		createdTask.ID(),
		200,
		map[string]string{
			"Test": "Test",
		},
		10,
	)

	Convey("TaskResultRepository INSERT and SELECT queries", t, func() {
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
