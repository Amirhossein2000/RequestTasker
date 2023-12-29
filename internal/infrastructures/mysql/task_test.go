package mysql

import (
	"context"
	"testing"

	"github.com/Amirhossein2000/RequestTasker/internal/domain/common"
	"github.com/Amirhossein2000/RequestTasker/internal/pkg/integration"
	"github.com/Amirhossein2000/RequestTasker/internal/pkg/test"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTaskRepository(t *testing.T) {
	conn, tearDown, err := integration.SetupMySQLContainer()
	if err != nil {
		t.Fatal(err)
	}
	defer tearDown()

	repo := NewTaskRepository(conn, common.TaskTable)

	task := test.NewTestTask()

	Convey("TaskRepository INSERT and SELECT queries", t, func() {
		Convey("Insert new task and check created task", func() {
			createdTask, err := repo.Create(context.Background(), task)
			So(err, ShouldBeNil)

			So(createdTask.ID(), ShouldNotEqual, int64(0))
			So(createdTask.Url(), ShouldEqual, task.Url())
			So(createdTask.Method(), ShouldEqual, task.Method())
			So(createdTask.Headers(), ShouldEqual, task.Headers())
			So(createdTask.Body(), ShouldEqual, task.Body())

			Convey("TaskRepository.GetByPublicID()", func() {
				foundTask, err := repo.GetByPublicID(context.Background(), createdTask.PublicID())
				So(err, ShouldBeNil)

				So(foundTask.PublicID(), ShouldEqual, createdTask.PublicID())
				So(foundTask.Url(), ShouldEqual, createdTask.Url())
				So(foundTask.Method(), ShouldEqual, createdTask.Method())
				So(foundTask.Headers(), ShouldEqual, createdTask.Headers())
				So(foundTask.Body(), ShouldEqual, createdTask.Body())
			})
		})
	})
}
