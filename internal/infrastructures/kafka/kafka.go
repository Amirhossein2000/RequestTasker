// main.go
package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Brokers           []string
	Topic             string
	GroupID           string
	Timeout           time.Duration
	NumPartitions     int
	ReplicationFactor int
}

type TaskEventRepository struct {
	config Config
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewTaskEventRepository(conf Config) *TaskEventRepository {
	dialer := &kafka.Dialer{
		Timeout:   conf.Timeout,
		DualStack: true,
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: conf.Brokers,
		Topic:   conf.Topic,
		Dialer:  dialer,
	})

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: conf.Brokers,
		GroupID: conf.GroupID,
		Topic:   conf.Topic,
		Dialer:  dialer,
	})

	return &TaskEventRepository{
		config: conf,
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

func (r *TaskEventRepository) CreateTopics(ctx context.Context) error {
	for _, broker := range r.config.Brokers {
		conn, err := kafka.DialContext(ctx, "tcp", broker)
		if err != nil {
			return err
		}

		err = conn.CreateTopics(kafka.TopicConfig{
			Topic:             r.config.Topic,
			NumPartitions:     r.config.NumPartitions,
			ReplicationFactor: r.config.ReplicationFactor,
		})
		if err != nil && err != kafka.TopicAlreadyExists {
			return err
		}
	}

	return nil
}
