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
	writer *kafka.Writer
}

func (p *PubTopic[MessageType]) Name() string {
	return p.writer.Topic
}

func (p *PubTopic[MessageType]) Pub(ctx context.Context, message MessageType) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Topic: p.writer.Topic,
		// TODO доделать эти поля
		Partition:     0,
		Offset:        0,
		HighWaterMark: 0,
		Key:           []byte(p.writer.Topic + " message"),
		Value:         message.Bytes(),
		// TODO возомжность настройки времени?
		Time: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("write message to kafka topic: %w", err)
	}
	return nil
}

func NewPubBroker[MessageType Message](writer *kafka.Writer) messagebroker.PubTopic[MessageType] {
	return &PubTopic[MessageType]{writer: writer}
}
