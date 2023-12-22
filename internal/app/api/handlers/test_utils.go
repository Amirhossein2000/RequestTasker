package handlers

import (
	"context"
	"net/http"
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
)

type testEnv struct {
	addr   string
	apiKey string
}

func setUpTestEnv() (*testEnv, func(), error) {
	addr, cleanup, err := integration.SetupKafkaContainer(context.Background())
	if err != nil {
		return nil, nil, err
	}

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
		return nil, nil, err
	}

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

	return &testEnv{
			addr:   testServerAddr,
			apiKey: apiKey,
		}, func() {
			testServer.Shutdown(context.Background())
			tearDown()
			cleanup()
		}, nil
}
