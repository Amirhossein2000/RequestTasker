// main.go
package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type TaskEventRepository struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewTaskEventRepository(brokers []string, topic string, groupID string) *TaskEventRepository {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
		Dialer:  dialer,
	})

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
		Dialer:  dialer,
	})

	return &TaskEventRepository{
		writer: writer,
		reader: reader,
	}
}

func (r *TaskEventRepository) Read(ctx context.Context) ([]byte, error) {
	msg, err := r.reader.ReadMessage(ctx)
	return msg.Value, err
}

func (r *TaskEventRepository) Write(ctx context.Context, value []byte) error {
	return r.writer.WriteMessages(ctx, kafka.Message{
		Value: value,
	})
}
