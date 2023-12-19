package kafka

import (
	"RequestTasker/internal/domian/dto"
	"RequestTasker/internal/pkg/integration"
	"context"
	"math/rand"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTaskEventRepository(t *testing.T) {
	Convey("TaskEventRepository Read and Write", t, func() {
		addr, cleanup, err := integration.SetupKafkaContainer(context.Background())
		So(err, ShouldBeNil)
		defer cleanup()

		repo := NewTaskEventRepository(
			KafkaConfig{
				Brokers:           []string{addr},
				Topic:             "test-topic",
				GroupID:           "test-group-id",
				Timeout:           time.Second * 10,
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
		)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err = repo.CreateTopics(ctx)
		So(err, ShouldBeNil)

		event := dto.TaskEvent{
			ID: rand.Int63(),
		}
		writeValue, err := event.Serialize()
		So(err, ShouldBeNil)

		err = repo.Write(ctx, writeValue)
		So(err, ShouldBeNil)
		readValue, err := repo.Read(ctx)
		So(err, ShouldBeNil)

		So(readValue, ShouldEqual, writeValue)
		receivedEvent, err := dto.NewTaskEvent(readValue)
		So(err, ShouldBeNil)
		So(event, ShouldEqual, *receivedEvent)
	})
}
