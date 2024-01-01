package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

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
	taskEventRepo := kafka.NewTaskEventRepository(*kafkaConf)

	taskRepo := mysql.NewTaskRepository(conn, "tasks")
	taskStatusRepo := mysql.NewTaskStatusRepository(conn, "task_statuses")
	taskResultRepo := mysql.NewTaskResultRepository(conn, "task_results")

	taskerService := tasker.NewTasker(
		taskEventRepo,
		taskRepo,
		taskStatusRepo,
		taskResultRepo,
		http.DefaultClient,
	)
	taskerService.Start(ctx)

	createTaskUsecase := usecases.NewCreateTaskUseCase(logger.Logger{}, taskRepo, taskStatusRepo, taskerService)
	getTaskUsecase := usecases.NewGetTaskUseCase(logger.Logger{}, taskRepo, taskStatusRepo, taskResultRepo)

	handler := handlers.NewHandler(os.Getenv("APP_API_KEY"), createTaskUsecase, getTaskUsecase)

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

func getMySQLConfig() (*mysql.Config, error) {
	addr := os.Getenv("MYSQL_ADDR")
	if addr == "" {
		return nil, errors.New("MYSQL_ADDR is not set")
	}

	user := os.Getenv("MYSQL_USER")
	if user == "" {
		return nil, errors.New("MYSQL_USER is not set")
	}

	password := os.Getenv("MYSQL_PASSWORD")
	if password == "" {
		return nil, errors.New("MYSQL_PASSWORD is not set")
	}

	database := os.Getenv("MYSQL_DB")
	if database == "" {
		return nil, errors.New("MYSQL_DB is not set")
	}

	return &mysql.Config{
		Addr:     addr,
		User:     user,
		Password: password,
		Database: database,
	}, nil
}

func getKafkaConfig() (*kafka.Config, error) {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		return nil, errors.New("KAFKA_BROKERS is not set")
	}

	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		return nil, errors.New("KAFKA_TOPIC is not set")
	}

	groupID := os.Getenv("KAFKA_GROUP_ID")
	if groupID == "" {
		return nil, errors.New("KAFKA_GROUP_ID is not set")
	}

	timeout, err := strconv.Atoi(os.Getenv("KAFKA_TIMEOUT_MILLISECONDS"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse KAFKA_TIMEOUT_MILLISECONDS: %v", err)
	}

	numPartitions, err := strconv.Atoi(os.Getenv("KAFKA_NUM_PARTITIONS"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse KAFKA_NUM_PARTITIONS: %v", err)
	}

	replicationFactor, err := strconv.Atoi(os.Getenv("KAFKA_REPLICATION_FACTOR"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse KAFKA_REPLICATION_FACTOR: %v", err)
	}

	return &kafka.Config{
		Brokers:           strings.Split(brokers, ","),
		Topic:             topic,
		GroupID:           groupID,
		Timeout:           time.Millisecond * time.Duration(timeout),
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}, nil
}
