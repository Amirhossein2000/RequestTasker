package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/Amirhossein2000/RequestTasker/internal/app/api"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/common"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/entities"
	"github.com/Amirhossein2000/RequestTasker/internal/pkg/test"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetTaskId(t *testing.T) {
	env, terminate, err := setUpTestEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer terminate()

	route := "/task/%s"
	ctx := context.Background()

	Convey("GetTaskId all possible responses", t, func() {
		Convey("return 400 when id is invalid", func() {
			id := "invalid_uuid"
			req, err := env.newReq(http.MethodGet, fmt.Sprintf(route, id), nil)
			So(err, ShouldBeNil)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)

			So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		})

		Convey("return 404 when task is not found", func() {
			body := api.GetTaskIdRequestObject{
				Id: uuid.NewString(),
			}

			req, err := env.newReq(http.MethodGet, fmt.Sprintf(route, body.Id), body)
			So(err, ShouldBeNil)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)

			So(resp.StatusCode, ShouldEqual, http.StatusNotFound)
		})

		Convey("return 200 when task is In_Progress", func() {
			task := test.NewTestTask()
			createdTask, err := env.taskRepository.Create(ctx, task)
			So(err, ShouldBeNil)

			status := entities.NewTaskStatus(createdTask.ID(), common.StatusIN_PROGRESS)
			_, err = env.taskStatusRepository.Create(ctx, status)
			So(err, ShouldBeNil)

			body := api.GetTaskIdRequestObject{
				Id: task.PublicID().String(),
			}

			req, err := env.newReq(http.MethodGet, fmt.Sprintf(route, body.Id), body)
			So(err, ShouldBeNil)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)

			So(resp.StatusCode, ShouldEqual, http.StatusOK)

			byteBody, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			respBody := api.GetTaskId200JSONResponse{}
			err = json.Unmarshal(byteBody, &respBody)
			So(err, ShouldBeNil)

			So(respBody.Id, ShouldEqual, createdTask.PublicID().String())
			So(string(respBody.Status), ShouldEqual, status.Status())
		})

		Convey("return 200 when task is Done", func() {
			task := test.NewTestTask()
			createdTask, err := env.taskRepository.Create(ctx, task)
			So(err, ShouldBeNil)

			status := entities.NewTaskStatus(createdTask.ID(), common.StatusDONE)
			_, err = env.taskStatusRepository.Create(ctx, status)
			So(err, ShouldBeNil)

			result := entities.NewTaskResult(
				createdTask.ID(),
				200,
				map[string]string{"test": "test"},
				10,
			)

			_, err = env.taskResultRepository.Create(ctx, result)
			So(err, ShouldBeNil)

			body := api.GetTaskIdRequestObject{
				Id: task.PublicID().String(),
			}
			req, err := env.newReq(http.MethodGet, fmt.Sprintf(route, body.Id), body)
			So(err, ShouldBeNil)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)

			So(resp.StatusCode, ShouldEqual, http.StatusOK)

			byteBody, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			respBody := api.GetTaskId200JSONResponse{}
			err = json.Unmarshal(byteBody, &respBody)
			So(err, ShouldBeNil)

			So(respBody.Id, ShouldEqual, createdTask.PublicID().String())
			So(string(respBody.Status), ShouldEqual, status.Status())
			So(*respBody.HttpStatusCode, ShouldEqual, result.StatusCode())
			So(*respBody.Headers, ShouldEqual, convertHeadersForResponse(result.Headers()))
			So(*respBody.Length, ShouldResemble, result.Length())
		})
	})
}
