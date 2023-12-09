package kafkapubsub

import (
	"context"
	"fmt"
	"github.com/achillescres/pkg/messagebroker"
	"github.com/segmentio/kafka-go"
	"time"
)

// PubTopic is kafka publish topic implementation.
// It uses kafka.Writer from kafka-go
// Specify your own MessageType with Message interface
type PubTopic[MessageType Message] struct {
	writer    *kafka.Writer
	partition int
	offset    int64
	timeFunc  func() time.Time
}

func NewPubTopic[MessageType Message](writer *kafka.Writer, partition int, offset int64, timeFunc func() time.Time) messagebroker.PubTopic[MessageType] {
	return &PubTopic[MessageType]{writer: writer, partition: partition, offset: offset, timeFunc: timeFunc}
}

func (p *PubTopic[MessageType]) Name() string {
	return fmt.Sprintf("%s/%d, offset=%d", p.writer.Topic, p.partition, p.offset)
}

func (p *PubTopic[MessageType]) Pub(ctx context.Context, message MessageType) error {
	rawMes, err := message.Bytes()
	if err != nil {
		return fmt.Errorf("get bytes from message: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Topic:     p.writer.Topic,
		Partition: p.partition,
		Offset:    p.offset,
		Key:       []byte(message.Key()),
		Value:     rawMes,
		Time:      p.timeFunc(),
	})
	if err != nil {
		return fmt.Errorf("write message to partition: %w", err)
	}
	return nil
}
