package main

import (
	"context"
	"net/http"
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
	ctx := context.Background()

	conn, err := mysql.NewMYSQLConn(mysql.MYSQLConfig{
		Addr:     "",
		User:     "",
		Password: "",
		Database: "",
	})
	if err != nil {
		panic(err)
	}

	taskRepo := mysql.NewTaskRepository(conn, "tasks")
	taskStatusRepo := mysql.NewTaskStatusRepository(conn, "taskStatuses")
	taskResultRepo := mysql.NewTaskResultRepository(conn, "taskResults")
	taskEventRepo := kafka.NewTaskEventRepository(kafka.KafkaConfig{
		Brokers:           []string{},
		Topic:             "",
		GroupID:           "",
		Timeout:           0,
		NumPartitions:     0,
		ReplicationFactor: 0,
	})

	taskerService := tasker.NewTasker(
		taskEventRepo,
		taskRepo,
		taskStatusRepo,
		taskResultRepo,
		http.DefaultClient,
	)
	go func() {
		err := taskerService.Start(ctx)
		if err != nil {
			panic(err)
		}
	}()

	createTaskUsecase := usecases.NewCreateTaskUseCase(logger.Logger{}, taskRepo, taskStatusRepo, taskerService)
	getTaskUsecase := usecases.NewGetTaskUseCase(logger.Logger{}, taskRepo, taskStatusRepo, taskResultRepo)

	handler := handlers.NewHandler("", createTaskUsecase, getTaskUsecase)

	httpServer := server.NewServer("", handler, []api.StrictMiddlewareFunc{handler.AuthMiddleware})
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
