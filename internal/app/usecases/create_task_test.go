package usecases

import (
	"errors"
	"math/rand"
	"testing"

	"RequestTasker/internal/app/services/logger"
	"RequestTasker/internal/domian/common"
	"RequestTasker/internal/domian/entities"
	"RequestTasker/internal/mocks"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateTaskUseCase_Execute(t *testing.T) {
	Convey("CreateTaskUseCase.Execute()", t, func() {
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
		createdTask := entities.BuildTask(
			rand.Int63(),
			newTask.CreatedAt(),
			newTask.PublicID(),
			newTask.Url(),
			newTask.Method(),
			newTask.Headers(),
			newTask.Body(),
		)
		status := entities.NewTaskStatus(createdTask.ID(), common.StatusNEW)

		Convey("When taskRepository.Create() returns error", func() {
			taskRepository.On("Create", newTask).Return(entities.Task{}, expectedErr)

			_, err := createTaskUseCase.Execute(newTask)
			So(err, ShouldEqual, common.InternalError)
		})

		Convey("When taskRepository.Create() works", func() {
			taskRepository.On("Create", newTask).Return(createdTask, nil)

			Convey("When taskStatusRepository.Create() returns error", func() {
				taskStatusRepository.On("Create", status).Return(expectedErr)

				_, err := createTaskUseCase.Execute(newTask)
				So(err, ShouldEqual, common.InternalError)
			})

			Convey("When taskStatusRepository.Create() works", func() {
				taskStatusRepository.On("Create", status).Return(nil)

				Convey("When requestTasker.RegisterTask() returns error", func() {
					requestTasker.On("RegisterTask", createdTask).Return(expectedErr)

					_, err := createTaskUseCase.Execute(newTask)
					So(err, ShouldEqual, common.InternalError)
				})

				Convey("When requestTasker.RegisterTask() works", func() {
					requestTasker.On("RegisterTask", createdTask).Return(nil)
	
					publicID, err := createTaskUseCase.Execute(newTask)
					So(err, ShouldBeNil)
					So(publicID, ShouldEqual, createdTask.PublicID())
				})
			})
		})
	})
}
