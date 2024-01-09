package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Amirhossein2000/RequestTasker/internal/app/api"
	"github.com/Amirhossein2000/RequestTasker/internal/app/api/handlers"
	"github.com/Amirhossein2000/RequestTasker/internal/app/services/logger"
	"github.com/Amirhossein2000/RequestTasker/internal/app/services/server"
	"github.com/Amirhossein2000/RequestTasker/internal/app/services/tasker"
	"github.com/Amirhossein2000/RequestTasker/internal/app/usecases"
	"github.com/Amirhossein2000/RequestTasker/internal/infrastructures/kafka"
	"github.com/Amirhossein2000/RequestTasker/internal/infrastructures/mysql"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mysqlConf, err := getMySQLConfig()
	if err != nil {
		panic(err)
	}
	conn, err := mysql.NewMYSQLConn(*mysqlConf)
	if err != nil {
		panic(err)
	}

	kafkaConf, err := getKafkaConfig()
	if err != nil {
		panic(err)
	}
	taskEventRepo, err := kafka.NewTaskEventRepository(ctx, *kafkaConf)
	if err != nil {
		panic(err)
	}

	logger, err := logger.NewLogger(false)
	if err != nil {
		panic(err)
	}

	taskRepo := mysql.NewTaskRepository(conn, "tasks")
	taskStatusRepo := mysql.NewTaskStatusRepository(conn, "task_statuses")
	taskResultRepo := mysql.NewTaskResultRepository(conn, "task_results")

	httpClient, err := getHttpClient()
	if err != nil {
		panic(err)
	}

	taskerService := tasker.NewTasker(
		logger,
		taskEventRepo,
		taskRepo,
		taskStatusRepo,
		taskResultRepo,
		httpClient,
	)
	taskerService.Start(ctx)

	createTaskUsecase := usecases.NewCreateTaskUseCase(taskRepo, taskStatusRepo, taskerService)
	getTaskUsecase := usecases.NewGetTaskUseCase(taskRepo, taskStatusRepo, taskResultRepo)

	handler := handlers.NewHandler(logger, os.Getenv("APP_API_KEY"), createTaskUsecase, getTaskUsecase)
	httpServer := server.NewServer(os.Getenv("APP_SERVER_ADDR"), handler, []api.StrictMiddlewareFunc{handler.AuthMiddleware})
	go func() {
		err := httpServer.Start()
		if err != nil {
			panic(err)
		}
	}()
	defer func() {
		err := httpServer.Shutdown(ctx)
		if err != nil {
			panic(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	os.Exit(0)
}
