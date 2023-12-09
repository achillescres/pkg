package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"os"
	"time"
)

type ReaderConfig struct {
	Brokers         []string
	GroupID         string
	Topic           string
	WaitingDuration time.Duration

	UseSASL            bool
	Username, Password string

	UseCA     bool
	CaAbsPath string
}

func NewReader(ctx context.Context, rc *ReaderConfig) (*kafka.Reader, error) {
	config := newReaderFromReaderConfig(*rc)

	dialer := &kafka.Dialer{
		Timeout:   rc.WaitingDuration,
		DualStack: true,
	}

	if rc.UseSASL {
		mechanism, err := scram.Mechanism(scram.SHA512, rc.Username, rc.Password)
		if err != nil {
			err := fmt.Errorf("error with scram:%s", err)
			return nil, err
		}
		dialer.SASLMechanism = mechanism
	}

	if rc.UseCA {
		rootCertPool := x509.NewCertPool()
		pem, err := os.ReadFile(rc.CaAbsPath)
		if err != nil {
			return nil, err
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, fmt.Errorf("error couldn't append certs to cert pool")
		}

		dialer.TLS = &tls.Config{
			RootCAs:            rootCertPool,
			InsecureSkipVerify: true,
		}
	}
	config.Dialer = dialer

	for _, addr := range config.Brokers {
		conn, err := dialer.DialContext(ctx, "tcp", addr)
		if err != nil {
			return nil, fmt.Errorf("failed to dial with broker %s: %s", addr, err.Error())
		}
		_, err = conn.ReadPartitions()
		if err != nil {
			return nil, fmt.Errorf("kafka read partitions failed: %s", err)
		}
	}

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
