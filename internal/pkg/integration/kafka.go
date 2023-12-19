package integration

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type testKafka struct {
	broker     *kafka.Conn
	zkAddr     string
	inUseCount int64
}

var tk *testKafka

const starterScript = "/usr/sbin/testcontainers_start.sh"
const kafkaLocalPort = "9092"

func SetupKafkaContainer() (string, *kafka.Conn, string, func(), error) {
	if tk != nil {
		atomic.AddInt64(&tk.inUseCount, 1)
		return "", tk.broker, tk.zkAddr,
			func() {
				atomic.AddInt64(&tk.inUseCount, -1)
			},
			nil
	}

	ctx := context.Background()

	// Start ZooKeeper container
	zkReq := testcontainers.ContainerRequest{
		Image:        "wurstmeister/zookeeper:latest",
		ExposedPorts: []string{"2181/tcp"},
		Env: map[string]string{
			"ALLOW_ANONYMOUS_LOGIN": "yes",
		},
		WaitingFor: wait.ForLog("binding to port 0.0.0.0/0.0.0.0:2181").WithStartupTimeout(30 * time.Second),
	}

	zkContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: zkReq,
		Started:          true,
	})
	if err != nil {
		return "", nil, "", nil, err
	}

	zkHost, err := zkContainer.Host(ctx)
	if err != nil {
		return "", nil, "", nil, err
	}

	zkPort, err := zkContainer.MappedPort(ctx, "2181")
	if err != nil {
		return "", nil, "", nil, err
	}

	zkAddr := fmt.Sprintf("%s:%s", zkHost, zkPort.Port())

	// Start Kafka container with ZooKeeper connection
	req := testcontainers.ContainerRequest{
		Image:        "confluentinc/cp-kafka:latest",
		ExposedPorts: []string{kafkaLocalPort + "/tcp"},
		Env: map[string]string{
			"KAFKA_BROKER_ID":                      "1",
			"KAFKA_LISTENERS":                      "PLAINTEXT://0.0.0.0:9092,BROKER://0.0.0.0:9093",
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP": "BROKER:PLAINTEXT,PLAINTEXT:PLAINTEXT",
			"KAFKA_INTER_BROKER_LISTENER_NAME":     "BROKER",
			"KAFKA_ZOOKEEPER_CONNECT":              zkAddr,
			"ALLOW_PLAINTEXT_LISTENER":             "yes",
			"KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE":  "true",
			"KAFKA_ENABLE_KRAFT":                   "no",
			"KAFKA_ADVERTISED_LISTENERS":           "PLAINTEXT://0.0.0.0:9092",
		},
	}


	
	kafkaC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", nil, "", nil, err
	}

	host, err := kafkaC.Host(ctx)
	if err != nil {
		return "", nil, "", nil, err
	}

	port, err := kafkaC.MappedPort(ctx, "9092")
	if err != nil {
		return "", nil, "", nil, err
	}

	brokerAddress := fmt.Sprintf("%s:%s", host, port.Port())
	conn, err := kafka.DialContext(ctx, "tcp", brokerAddress)
	if err != nil {
		return "", nil, "", nil, err
	}

	tk = &testKafka{
		broker:     conn,
		zkAddr:     zkAddr,
		inUseCount: 1,
	}

	time.Sleep(time.Second * 30)

	return brokerAddress, conn, zkAddr,
		func() {
			inUseCount := atomic.AddInt64(&tk.inUseCount, -1)
			if inUseCount < 1 {
				kafkaC.Terminate(ctx)
				zkContainer.Terminate(ctx)
				conn.Close()
			}
		},
		nil
}
