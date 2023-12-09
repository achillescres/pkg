package kafka

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"time"
)

type WriterConfig struct {
	// The list of brokers used to discover the partitions available on the
	// kafka cluster.
	//
	// This field is required, attempting to create a writer with an empty list
	// of brokers will panic.
	Brokers []string

	// The topic that the writer will produce messages to.
	//
	// If provided, this will be used to set the topic for all produced messages.
	// If not provided, each Message must specify a topic for itself. This must be
	// mutually exclusive, otherwise the Writer will return an error.
	Topic string

	// The dialer used by the writer to establish connections to the kafka
	// cluster.
	//
	// If nil, the default dialer is used instead.
	Dialer *kafka.Dialer

	// Limit on how many attempts will be made to deliver a message.
	//
	// The default is to try at most 10 times.
	MaxAttempts *int

	// Limit on how many messages will be buffered before being sent to a
	// partition.
	//
	// The default is to use a target batch size of 100 messages.
	BatchSize *int

	// Limit the maximum size of a request in bytes before being sent to
	// a partition.
	//
	// The default is to use a kafka default value of 1048576.
	BatchBytes *int

	// Time limit on how often incomplete message batches will be flushed to
	// kafka.
	//
	// The default is to flush at least every second.
	BatchTimeout *time.Duration

	// Timeout for read operations performed by the Writer.
	//
	// Defaults to 10 seconds.
	ReadTimeout *time.Duration

	// Timeout for write operation performed by the Writer.
	//
	// Defaults to 10 seconds.
	WriteTimeout *time.Duration

	// Number of acknowledges from partition replicas required before receiving
	// a response to a produce request. The default is -1, which means to wait for
	// all replicas, and a value above 0 is required to indicate how many replicas
	// should acknowledge a message to be considered successful.
	RequiredAcks *int

	// Setting this flag to true causes the WriteMessages method to never block.
	// It also means that errors are ignored since the caller will not receive
	// the returned value. Use this only if you don't care about guarantees of
	// whether the messages were written to kafka.
	Async bool

	// If not nil, specifies a logger used to report internal changes within the
	// writer.
	Logger kafka.Logger

	// ErrorLogger is the logger used to report errors. If nil, the writer falls
	// back to using Logger instead.
	ErrorLogger kafka.Logger
}

func NewWriter(ctx context.Context, config WriterConfig) (*kafka.Writer, error) {
	if config.Dialer == nil {
		return nil, errors.New("dialer is nil")
	}
	if config.MaxAttempts == nil {
		config.MaxAttempts = new(int)
		*config.MaxAttempts = 10
	}
	if config.BatchSize == nil {
		config.BatchSize = new(int)
		*config.BatchSize = 100
	}
	if config.BatchBytes == nil {
		config.BatchBytes = new(int)
		*config.BatchBytes = 1048576
	}
	if config.BatchTimeout == nil {
		config.BatchTimeout = new(time.Duration)
		*config.BatchTimeout = time.Second + time.Millisecond*500
	}
	if config.ReadTimeout == nil {
		config.ReadTimeout = new(time.Duration)
		*config.ReadTimeout = time.Second * 10
	}
	if config.WriteTimeout == nil {
		config.WriteTimeout = new(time.Duration)
		*config.WriteTimeout = time.Second * 10
	}
	if config.RequiredAcks == nil {
		config.RequiredAcks = new(int)
		*config.RequiredAcks = -1
	}
	if config.Logger == nil {
		return nil, errors.New("logger is nil")
	}
	if config.ErrorLogger == nil {
		return nil, errors.New("error logger is nil")
	}
	// TODO WARNING need to keep kafka-go version lower than v1.0
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      config.Brokers,
		Topic:        config.Topic,
		Dialer:       config.Dialer,
		MaxAttempts:  *config.MaxAttempts,
		BatchSize:    *config.BatchSize,
		BatchBytes:   *config.BatchBytes,
		BatchTimeout: *config.BatchTimeout,
		ReadTimeout:  *config.ReadTimeout,
		WriteTimeout: *config.WriteTimeout,
		RequiredAcks: *config.RequiredAcks,
		Async:        config.Async,
		Logger:       config.Logger,
		ErrorLogger:  config.ErrorLogger,
	})
	return w, nil
}
