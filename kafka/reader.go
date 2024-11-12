package kafka

import (
	"github.com/segmentio/kafka-go"
	"time"
)

type ReaderConfig struct {
	Brokers []string
	Dialer  *kafka.Dialer
	GroupID string
	Topic   string
	Timeout time.Duration
}

func NewReader(rc ReaderConfig) (*kafka.Reader, error) {
	config := newReaderFromReaderConfig(rc)

	r := kafka.NewReader(*config)

	return r, nil
}

func newReaderFromReaderConfig(rc ReaderConfig) *kafka.ReaderConfig {
	conf := &kafka.ReaderConfig{
		Dialer:  rc.Dialer,
		Brokers: rc.Brokers,
		GroupID: rc.GroupID,
		Topic:   rc.Topic,
		MaxWait: rc.Timeout,
	}
	return conf
}
