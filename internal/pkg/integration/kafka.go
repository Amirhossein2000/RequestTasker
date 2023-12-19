package integration

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

func SetupKafkaContainer(ctx context.Context) (string, func(), error) {
	kafkaC, err := kafka.RunContainer(ctx,
		kafka.WithClusterID("test-cluster"),
	)
	if err != nil {
		panic(err)
	}
	
	host, err := kafkaC.Host(ctx)
	if err != nil {
		panic(err)
	}

	port, err := kafkaC.MappedPort(ctx, "9093")
	if err != nil {
		panic(err)
	}

	brokerAddress := fmt.Sprintf("%s:%s", host, port.Port())

	return brokerAddress,
		func() {
			if err := kafkaC.Terminate(ctx); err != nil {
				panic(err)
			}
		},
		err
}
