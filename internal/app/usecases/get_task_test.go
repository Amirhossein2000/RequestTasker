package usecases

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"RequestTasker/internal/app/services/logger"
	"RequestTasker/internal/domain/common"
	"RequestTasker/internal/domain/entities"
	"RequestTasker/internal/mocks"
	"RequestTasker/internal/pkg/test"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetTaskUseCase_Execute(t *testing.T) {
	Convey("GetTaskUseCase.Execute()", t, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		logger := logger.NewLogger()
		taskRepository := mocks.NewTaskRepositoryMock(t)
		taskStatusRepository := mocks.NewTaskStatusRepositoryMock(t)
		taskResultRepository := mocks.NewTaskResultRepositoryMock(t)
		getTaskUseCase := NewGetTaskUseCase(
			logger,
			taskRepository,
			taskStatusRepository,
			taskResultRepository,
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
		expectedStatus := entities.NewTaskStatus(expectedTask.ID(), common.StatusNEW)

		Convey("When taskRepository.GetByPublicID() returns error", func() {
			taskRepository.
				On("GetByPublicID", ctx, expectedTask.PublicID()).
				Return(nil, expectedErr)

			_, _, _, err := getTaskUseCase.Execute(ctx, expectedTask.PublicID())
			So(err, ShouldEqual, common.ErrInternal)
		})

		Convey("When taskRepository.GetByPublicID() works", func() {
			taskRepository.
				On("GetByPublicID", ctx, expectedTask.PublicID()).
				Return(&expectedTask, nil)

			Convey("When taskStatusRepository.GetLatestByTaskID() returns error", func() {
				taskStatusRepository.
					On("GetLatestByTaskID", ctx, expectedTask.ID()).
					Return(&expectedStatus, expectedErr)

				_, _, _, err := getTaskUseCase.Execute(ctx, expectedTask.PublicID())
				So(err, ShouldEqual, common.ErrInternal)
			})

			Convey("When taskStatusRepository.GetLatestByTaskID() returns status without result", func() {
				taskStatusRepository.
					On("GetLatestByTaskID", ctx, expectedTask.ID()).
					Return(&expectedStatus, nil)

				task, status, _, err := getTaskUseCase.Execute(ctx, expectedTask.PublicID())
				So(err, ShouldBeNil)
				So(task, ShouldResemble, &expectedTask)
				So(status, ShouldResemble, &expectedStatus)
			})

			Convey("When taskStatusRepository.GetLatestByTaskID() returns status with result", func() {
				expectedStatus := entities.NewTaskStatus(expectedTask.ID(), common.StatusDONE)
				expectedResult := entities.NewTaskResult(
					expectedTask.ID(),
					200,
					nil,
					10,
				)
				taskStatusRepository.
					On("GetLatestByTaskID", ctx, expectedTask.ID()).
					Return(&expectedStatus, nil)

				Convey("When taskResultRepository.GetByTaskID() returns error", func() {
					taskResultRepository.
						On("GetByTaskID", ctx, expectedTask.ID()).
						Return(nil, expectedErr)

					_, _, _, err := getTaskUseCase.Execute(ctx, expectedTask.PublicID())
					So(err, ShouldEqual, common.ErrInternal)
				})

				Convey("When taskResultRepository.GetByTaskID() works", func() {
					taskResultRepository.
						On("GetByTaskID", ctx, expectedTask.ID()).
						Return(&expectedResult, nil)

					task, status, result, err := getTaskUseCase.Execute(ctx, expectedTask.PublicID())
					So(err, ShouldBeNil)
					So(task, ShouldResemble, &expectedTask)
					So(status, ShouldResemble, &expectedStatus)
					So(result, ShouldResemble, &expectedResult)
				})
			})
		})
	})
}
