package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Amirhossein2000/RequestTasker/internal/app/api"
	"github.com/Amirhossein2000/RequestTasker/internal/app/services/logger"
	"github.com/Amirhossein2000/RequestTasker/internal/app/services/server"
	"github.com/Amirhossein2000/RequestTasker/internal/app/services/tasker"
	"github.com/Amirhossein2000/RequestTasker/internal/app/usecases"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/common"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/entities"
	"github.com/Amirhossein2000/RequestTasker/internal/infrastructures/kafka"
	"github.com/Amirhossein2000/RequestTasker/internal/infrastructures/mysql"
	"github.com/Amirhossein2000/RequestTasker/internal/pkg/integration"
	"github.com/labstack/gommon/random"
)

type testEnv struct {
	addr   string
	apiKey string

	taskRepository       entities.TaskRepository
	taskStatusRepository entities.TaskStatusRepository
	taskResultRepository entities.TaskResultRepository
}

func (e *testEnv) newReq(method string, route string, body any) (*http.Request, error) {
	rBody := bytes.NewBuffer(nil)

	if body != nil {
		byteBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		rBody = bytes.NewBuffer(byteBody)
		if err != nil {
			return nil, err
		}
	}

	u, err := url.JoinPath(fmt.Sprintf("http://%s", e.addr), route)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u, rBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", e.apiKey)

	return req, nil
}

func setUpTestEnv() (*testEnv, func(), error) {
	addr, cleanup, err := integration.SetupKafkaContainer(context.Background())
	if err != nil {
		return nil, nil, err
	}

	taskEventRepository := kafka.NewTaskEventRepository(
		kafka.Config{
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
		return nil, nil, err
	}

	taskRepository := mysql.NewTaskRepository(conn, common.TaskTable)
	taskStatusRepository := mysql.NewTaskStatusRepository(conn, common.TaskStatusTable)
	taskResultRepository := mysql.NewTaskResultRepository(conn, common.TaskResultTable)

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

	return &testEnv{
			addr:                 testServerAddr,
			apiKey:               apiKey,
			taskRepository:       taskRepository,
			taskStatusRepository: taskStatusRepository,
			taskResultRepository: taskResultRepository,
		}, func() {
			testServer.Shutdown(context.Background())
			tearDown()
			cleanup()
		}, nil
}
