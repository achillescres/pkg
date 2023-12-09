package kafkapubsub

import (
	"context"
	"github.com/achillescres/pkg/messagebroker"
	"github.com/segmentio/kafka-go"
)

// SubTopic is kafka subscribe topic implementation.
// It uses kafka.Reader from kafka-go
// Specify your own MessageType with Message interface
type SubTopic[MessageType Message] struct {
	reader *kafka.Reader
}

func NewSubTopic[MessageType Message](reader *kafka.Reader) messagebroker.SubTopic[MessageType] {
	return &SubTopic[MessageType]{reader: reader}
}

func (s *SubTopic[MessageType]) Name() string {
	return s.reader.Stats().Topic
}

func (s *SubTopic[MessageType]) Sub(callback messagebroker.Callback[MessageType]) (messagebroker.CancelSubscription, error) {
	// TODO ctx откуда?
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			rawMes, err := s.reader.ReadMessage(ctx)
			if err != nil {
				// TODO add log
				return
			}

			var mes MessageType
			err = mes.Scan(rawMes.Value)
			if err != nil {
				// TODO add log REALLY ADD LOGS
				return
			}

			// TODO maybe add Goard panic security
			callback(mes)
		}
	}()

	return messagebroker.CancelSubscription(cancel), nil
}
