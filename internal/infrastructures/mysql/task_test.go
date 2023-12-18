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

func TestTaskRepository(t *testing.T) {
	Convey("TaskRepository INSERT and SELECT queries", t, func() {
		conn, tearDown, err := integration.SetupMySQLContainer()
		So(err, ShouldBeNil)
		defer tearDown()

		session := conn.NewSession(nil)
		repo := NewTaskRepository(session, "tasks")

		task := entities.NewTask(
			"https://example.com",
			"GET",
			map[string]interface{}{"Authorization": "Bearer token"},
			map[string]interface{}{"key": "value"},
		)

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

				So(createdTask.PublicID(), ShouldEqual, foundTask.PublicID())
				So(createdTask.Url(), ShouldEqual, foundTask.Url())
				So(createdTask.Method(), ShouldEqual, foundTask.Method())
				So(createdTask.Headers(), ShouldEqual, foundTask.Headers())
				So(createdTask.Body(), ShouldEqual, foundTask.Body())
			})
		})
	})
}
