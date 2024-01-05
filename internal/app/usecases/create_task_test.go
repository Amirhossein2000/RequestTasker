package usecases

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/Amirhossein2000/RequestTasker/internal/domain/common"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/entities"
	"github.com/Amirhossein2000/RequestTasker/internal/mocks"
	"github.com/Amirhossein2000/RequestTasker/internal/pkg/test"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateTaskUseCase_Execute(t *testing.T) {
	Convey("CreateTaskUseCase.Execute()", t, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		taskRepository := mocks.NewTaskRepositoryMock(t)
		taskStatusRepository := mocks.NewTaskStatusRepositoryMock(t)
		tasker := mocks.NewTaskerMock(t)
		createTaskUseCase := NewCreateTaskUseCase(
			taskRepository,
			taskStatusRepository,
			tasker,
		)

		expectedErr := errors.New("expectedErr")
		newTask := test.NewTestTask()
		expectedTask := entities.BuildTask(
			rand.Int63(),
			newTask.CreatedAt(),
			newTask.PublicID(),
			newTask.Url(),
			newTask.Method(),
			newTask.Headers(),
			newTask.Body(),
		)
		expectedStatus := entities.NewTaskStatus(expectedTask.ID(), common.StatusNew)

		Convey("When taskRepository.Create() returns error", func() {
			taskRepository.
				On("Create", ctx, newTask).
				Return(nil, expectedErr)

			_, err := createTaskUseCase.Execute(ctx, newTask)
			So(err, ShouldEqual, expectedErr)
		})

		Convey("When taskRepository.Create() works", func() {
			taskRepository.
				On("Create", ctx, newTask).
				Return(&expectedTask, nil)

			Convey("When taskStatusRepository.Create() returns error", func() {
				taskStatusRepository.
					On("Create", ctx, expectedStatus).
					Return(nil, expectedErr)

				_, err := createTaskUseCase.Execute(ctx, newTask)
				So(err, ShouldEqual, expectedErr)
			})

			Convey("When taskStatusRepository.Create() works", func() {
				taskStatusRepository.
					On("Create", ctx, expectedStatus).
					Return(&expectedStatus, nil)

				Convey("When tasker.RegisterTask() returns error", func() {
					tasker.
						On("RegisterTask", ctx, expectedTask).
						Return(expectedErr)

					_, err := createTaskUseCase.Execute(ctx, newTask)
					So(err, ShouldEqual, expectedErr)
				})

				Convey("When tasker.RegisterTask() works", func() {
					tasker.
						On("RegisterTask", ctx, expectedTask).
						Return(nil)

					publicID, err := createTaskUseCase.Execute(ctx, newTask)
					So(err, ShouldBeNil)
					So(publicID, ShouldEqual, expectedTask.PublicID())
				})
			})
		})
	})
}
