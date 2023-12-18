package usecases

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"RequestTasker/internal/app/services/logger"
	"RequestTasker/internal/domian/common"
	"RequestTasker/internal/domian/entities"
	"RequestTasker/internal/mocks"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateTaskUseCase_Execute(t *testing.T) {
	Convey("CreateTaskUseCase.Execute()", t, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		logger := logger.NewLogger()
		taskRepository := mocks.NewTaskRepositoryMock(t)
		taskStatusRepository := mocks.NewTaskStatusRepositoryMock(t)
		requestTasker := mocks.NewRequestTaskerMock(t)
		createTaskUseCase := NewCreateTaskUseCase(
			logger,
			taskRepository,
			taskStatusRepository,
			requestTasker,
		)

		expectedErr := errors.New("expectedErr")
		newTask := entities.NewTask("url", "method", nil, nil)
		expectedTask := entities.BuildTask(
			rand.Int63(),
			newTask.CreatedAt(),
			newTask.PublicID(),
			newTask.Url(),
			newTask.Method(),
			newTask.Headers(),
			newTask.Body(),
		)
		expectedStatus := entities.NewTaskStatus(expectedTask.ID(), common.StatusNEW)

		Convey("When taskRepository.Create() returns error", func() {
			taskRepository.
				On("Create", ctx, newTask).
				Return(nil, expectedErr)

			_, err := createTaskUseCase.Execute(ctx, newTask)
			So(err, ShouldEqual, common.InternalError)
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
				So(err, ShouldEqual, common.InternalError)
			})

			Convey("When taskStatusRepository.Create() works", func() {
				taskStatusRepository.
					On("Create", ctx, expectedStatus).
					Return(&expectedStatus, nil)

				Convey("When requestTasker.RegisterTask() returns error", func() {
					requestTasker.
						On("RegisterTask", ctx, expectedTask).
						Return(expectedErr)

					_, err := createTaskUseCase.Execute(ctx, newTask)
					So(err, ShouldEqual, common.InternalError)
				})

				Convey("When requestTasker.RegisterTask() works", func() {
					requestTasker.
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
