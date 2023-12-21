package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Amirhossein2000/RequestTasker/internal/app/api"
	"github.com/Amirhossein2000/RequestTasker/internal/app/services/logger"
	"github.com/Amirhossein2000/RequestTasker/internal/app/services/server"
	"github.com/Amirhossein2000/RequestTasker/internal/app/services/tasker"
	"github.com/Amirhossein2000/RequestTasker/internal/app/usecases"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/common"
	"github.com/Amirhossein2000/RequestTasker/internal/infrastructures/kafka"
	"github.com/Amirhossein2000/RequestTasker/internal/infrastructures/mysql"
	"github.com/Amirhossein2000/RequestTasker/internal/pkg/integration"
	"github.com/labstack/gommon/random"
	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPostTask(t *testing.T) {
	addr, cleanup, err := integration.SetupKafkaContainer(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	taskEventRepository := kafka.NewTaskEventRepository(
		kafka.KafkaConfig{
			Brokers:           []string{addr},
			Topic:             "test-topic",
			GroupID:           "test-group-id",
			Timeout:           time.Second * 10,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	)

	conn, tearDown, err := integration.SetupMySQLContainer()
	if err != nil {
		t.Fatal(err)
	}
	defer tearDown()

	taskRepository := mysql.NewTaskRepository(conn.NewSession(nil), common.TaskTable)
	taskStatusRepository := mysql.NewTaskStatusRepository(conn.NewSession(nil), common.TaskStatusTable)
	taskResultRepository := mysql.NewTaskResultRepository(conn.NewSession(nil), common.TaskResultTable)

	tasker := tasker.NewTasker(
		taskEventRepository,
		taskRepository,
		taskStatusRepository,
		taskResultRepository,
		http.DefaultClient,
	)

	createTaskUseCase := usecases.NewCreateTaskUseCase(
		logger.Logger{},
		taskRepository,
		taskStatusRepository,
		tasker,
	)

	getTaskUseCase := usecases.NewGetTaskUseCase(
		logger.Logger{},
		taskRepository,
		taskStatusRepository,
		taskResultRepository,
	)

	apiKey := random.String(32)

	testHandler := NewHandler(
		apiKey,
		createTaskUseCase,
		getTaskUseCase,
	)

	testServerAddr := "localhost:7777"
	testServer := server.NewServer(
		testServerAddr,
		testHandler,
		[]api.StrictMiddlewareFunc{
			testHandler.AuthMiddleware,
		},
	)
	go testServer.Start()
	defer testServer.Shutdown(context.Background())

	Convey("PostTask all possible responses", t, func() {
		Convey("return 400 when url is invalid", func() {
			body := api.PostTaskJSONRequestBody{
				Body: lo.ToPtr(`{"test":"test"}`),
				Headers: &map[string]interface{}{
					"test": "test",
				},
				Method: "GET",
				Url:    "InvalidURL",
			}

			byteBody, err := json.Marshal(body)
			So(err, ShouldBeNil)

			rBody := bytes.NewBuffer(byteBody)
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/task", testServerAddr), rBody)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", apiKey)
			So(err, ShouldBeNil)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)

			So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		})

		Convey("return 400 when method is invalid", func() {
			body := api.PostTaskJSONRequestBody{
				Body: lo.ToPtr(`{"test":"test"}`),
				Headers: &map[string]interface{}{
					"test": "test",
				},
				Method: "InvalidMethod",
				Url:    "www.google.com",
			}

			byteBody, err := json.Marshal(body)
			So(err, ShouldBeNil)

			rBody := bytes.NewBuffer(byteBody)
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/task", testServerAddr), rBody)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", apiKey)
			So(err, ShouldBeNil)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)

			So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		})

		Convey("return 401 when api-key is wrong", func() {
			body := api.PostTaskJSONRequestBody{
				Body: lo.ToPtr(`{"test":"test"}`),
				Headers: &map[string]interface{}{
					"test": "test",
				},
				Method: "GET",
				Url:    "www.google.com",
			}

			byteBody, err := json.Marshal(body)
			So(err, ShouldBeNil)

			rBody := bytes.NewBuffer(byteBody)
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/task", testServerAddr), rBody)
			req.Header.Set("Content-Type", "application/json")
			So(err, ShouldBeNil)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)

			So(resp.StatusCode, ShouldEqual, http.StatusUnauthorized)
		})

		Convey("return 201 when everything is fine", func() {
			body := api.PostTaskJSONRequestBody{
				Body: lo.ToPtr(`{"test":"test"}`),
				Headers: &map[string]interface{}{
					"test": "test",
				},
				Method: "GET",
				Url:    "https://www.google.com",
			}

			byteBody, err := json.Marshal(body)
			So(err, ShouldBeNil)

			rBody := bytes.NewBuffer(byteBody)
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/task", testServerAddr), rBody)
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", apiKey)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)

			So(resp.StatusCode, ShouldEqual, http.StatusCreated)
		})
	})
}
