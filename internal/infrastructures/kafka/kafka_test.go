package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/Amirhossein2000/RequestTasker/internal/domain/dto"
	"github.com/Amirhossein2000/RequestTasker/internal/pkg/integration"
	"github.com/google/uuid"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTaskEventRepository(t *testing.T) {
	addr, cleanup, err := integration.SetupKafkaContainer(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	repo := NewTaskEventRepository(
		Config{
			Brokers:           []string{addr},
			Topic:             "test-topic",
			GroupID:           "test-group-id",
			Timeout:           time.Second * 10,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	err = repo.CreateTopics(ctx)
	if err != nil {
		t.Fatal(err)
	}

	Convey("TaskEventRepository Read and Write", t, func() {
		for i := 0; i < 3; i++ {
			event := dto.TaskEvent{
				PublicID: uuid.New(),
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
		}
	})
}
