package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

type ReaderConfig struct {
	Dialer  *kafka.Dialer
	Brokers []string

	GroupID string
	Topic   string
	Timeout time.Duration

	UseSASL            bool
	Username, Password string

	UseCA     bool
	CaAbsPath string
}

func NewReader(ctx context.Context, rc ReaderConfig) (*kafka.Reader, error) {
	config := newReaderFromReaderConfig(rc)

	r := kafka.NewReader(*config)

	return r, nil
}

func newReaderFromReaderConfig(rc ReaderConfig) *kafka.ReaderConfig {
	conf := &kafka.ReaderConfig{
		Brokers: rc.Brokers,
		GroupID: rc.GroupID,
		Topic:   rc.Topic,
	}
	return conf
}
