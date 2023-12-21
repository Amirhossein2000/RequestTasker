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
	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPostTask(t *testing.T) {
	Convey("PostTask all possible responses", t, func() {
		addr, cleanup, err := integration.SetupKafkaContainer(context.Background())
		So(err, ShouldBeNil)
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

		session, tearDown, err := integration.SetupMySQLContainer()
		So(err, ShouldBeNil)
		defer tearDown()

		taskRepository := mysql.NewTaskRepository(session, common.TaskTable)
		taskStatusRepository := mysql.NewTaskStatusRepository(session, common.TaskStatusTable)
		taskResultRepository := mysql.NewTaskResultRepository(session, common.TaskResultTable)

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

		testHandler := NewHandler(
			"API-Key",
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
	})
}
