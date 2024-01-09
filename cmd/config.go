package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Amirhossein2000/RequestTasker/internal/infrastructures/kafka"
	"github.com/Amirhossein2000/RequestTasker/internal/infrastructures/mysql"
)

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

func getHttpClient() (*http.Client, error) {
	httpClient := http.DefaultClient
	clientTimeout, err := strconv.Atoi(os.Getenv("APP_HTTP_CLIENT_TIMEOUT_MILLISECONDS"))
	if err != nil {
		panic(err)
	}
	httpClient.Timeout = time.Millisecond * time.Duration(clientTimeout)

	return httpClient, nil
}
