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
	"github.com/Amirhossein2000/RequestTasker/internal/pkg/test"
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

var tE *testEnv

func setUpTestEnv() (*testEnv, func(), error) {
	if tE != nil {
		return tE, func() {}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	addr, cleanup, err := integration.SetupKafkaContainer(ctx)
	if err != nil {
		return nil, nil, err
	}

	taskEventRepository, err := kafka.NewTaskEventRepository(ctx,
		kafka.Config{
			Brokers:           []string{addr},
			Topic:             "test-topic",
			GroupID:           "test-group-id",
			Timeout:           time.Second * 10,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	)
	if err != nil {
		return nil, nil, err
	}

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

	logger, err := logger.NewLogger(true)
	if err != nil {
		return nil, nil, err
	}

	createTaskUseCase := usecases.NewCreateTaskUseCase(
		logger,
		taskRepository,
		taskStatusRepository,
		tasker,
	)

	getTaskUseCase := usecases.NewGetTaskUseCase(
		logger,
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

	port, err := test.GetAvailablePort()
	if err != nil {
		return nil, nil, err
	}

	testServerAddr := fmt.Sprintf("127.0.0.1:%d", port)
	testServer := server.NewServer(
		testServerAddr,
		testHandler,
		[]api.StrictMiddlewareFunc{
			testHandler.AuthMiddleware,
		},
	)
	go testServer.Start()

	tE = &testEnv{
		addr:                 testServerAddr,
		apiKey:               apiKey,
		taskRepository:       taskRepository,
		taskStatusRepository: taskStatusRepository,
		taskResultRepository: taskResultRepository,
	}

	return tE, func() {
		testServer.Shutdown(ctx)
		tearDown()
		cleanup()
	}, nil
}
