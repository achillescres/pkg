package kafkapubsub

import (
	"context"
	"fmt"
	"github.com/achillescres/pkg/infrastructure/tube"
	"github.com/achillescres/pkg/messagebroker"
	"github.com/segmentio/kafka-go"
)

type Commit int

const (
	AutoCommit = iota
	ManualCommit
)

// SubTopic is kafka subscribe topic implementation.
// It uses kafka.Reader from kafka-go
// Specify your own MessageType with Message interface
// Don't use it directly instead use NewSubTopic!
type SubTopic[MessageType Message[MessageType]] struct {
	reader  *kafka.Reader
	errTube tube.Error
	// commit determines whether Sub will commit message after callback invoke or immediately
	commit Commit
}

func NewSubTopic[MessageType Message[MessageType]](reader *kafka.Reader, errTube tube.Error, commit Commit) messagebroker.SubTopic[MessageType] {
	return &SubTopic[MessageType]{reader: reader, errTube: errTube, commit: commit}
}

func (s *SubTopic[MessageType]) Name() string {
	return s.reader.Stats().Topic
}

func (s *SubTopic[MessageType]) Sub(callback messagebroker.Callback[MessageType]) (messagebroker.CancelSubscription, error) {
	// TODO ctx откуда?
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			rawMes, err := s.reader.FetchMessage(ctx)
			if err != nil {
				s.errTube(fmt.Errorf("fetch message: %w", err))
				return
			}
			if s.commit == AutoCommit {
				err = s.reader.CommitMessages(ctx, rawMes)
				if err != nil {
					s.errTube(fmt.Errorf("commit message: %w", err))
					return
				}
			}

			var mes MessageType
			mes, err = mes.Unmarshal(rawMes.Value)
			if err != nil {
				s.errTube(fmt.Errorf("scan message's value: %w", err))
				return
			}

			// TODO maybe add Goard panic security
			callback(mes)

			if s.commit == ManualCommit {
				err = s.reader.CommitMessages(ctx, rawMes)
				if err != nil {
					s.errTube(fmt.Errorf("commit message: %w", err))
					return
				}
			}
		}
	}()

	return messagebroker.CancelSubscription(cancel), nil
}
