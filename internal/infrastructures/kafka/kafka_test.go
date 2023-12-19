package kafka

import (
	"RequestTasker/internal/pkg/integration"
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestTaskEventRepository(t *testing.T) {
	Convey("TaskEventRepository Read and Write", t, func() {
		addr, conn, _, cleanup, err := integration.SetupKafkaContainer()
		if err != nil {
			t.Fatalf("Failed to start Kafka container: %v", err)
		}
		defer cleanup()

		repo := NewTaskEventRepository([]string{addr}, "test-topic", "test-group-id")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		_, err = conn.Write([]byte("HI"))
		So(err, ShouldBeNil)

		err = repo.Write(ctx, []byte("{\"Test_key\": \"Test Value\"}"))
		So(err, ShouldBeNil)
	})
}
