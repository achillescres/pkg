package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/achillescres/pkg/utils"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"os"
	"time"
)

type ScramAlgorithm uint

func (sa ScramAlgorithm) Scram() scram.Algorithm {
	switch sa {
	case SHA256:
		return scram.SHA256
	case SHA512:
		return scram.SHA512
	default:
		panic("invalid ScramAlgorithm")
	}
}

func (sa ScramAlgorithm) Validate() error {
	switch sa {
	case SHA256:
		return nil
	case SHA512:
		return nil
	default:
		return fmt.Errorf("invalid ScramAlgorithm: %d", sa)
	}
}

const ( // shit but reverse compatible shit
	SHA256 ScramAlgorithm = iota + 1
	SHA512
)

type DialerConfig struct {
	Username, Password      string
	UseSASL                 bool
	UseCA                   bool
	CaAbsPath               string
	Timeout                 time.Duration
	VerifyBrokerCertificate bool
	// Default 0 is invalid
	ScramAlgorithm ScramAlgorithm
}

func (dc DialerConfig) validate() error {
	ew := utils.NewErrorWrapper("DialerConfig - validate")

	if err := dc.ScramAlgorithm.Validate(); err != nil {
		return ew(fmt.Errorf("validate ScramAlgorithm: %w", err))
	}

	return nil
}

func NewDialer(config DialerConfig) (*kafka.Dialer, error) {
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	dialer := kafka.Dialer{
		Timeout:   config.Timeout,
		DualStack: true,
	}

	if config.UseSASL {
		mechanism, err := scram.Mechanism(config.ScramAlgorithm.Scram(), config.Username, config.Password)
		if err != nil {
			return nil, fmt.Errorf("create mechanism: %w", err)
		}
		dialer.SASLMechanism = mechanism
	}

	tlsc := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	if config.UseCA {
		rootCertPool := x509.NewCertPool()
		pem, err := os.ReadFile(config.CaAbsPath)
		if err != nil {
			return nil, fmt.Errorf("read ca path: %w", err)
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, fmt.Errorf("append certs from file to cert pool")
		}
		tlsc.RootCAs = rootCertPool
		tlsc.InsecureSkipVerify = !config.VerifyBrokerCertificate
	}

	dialer.TLS = tlsc

	return &dialer, nil
}

func TestDialer(ctx context.Context, d *kafka.Dialer, brokers []string) error {
	for _, addr := range brokers {
		conn, err := d.DialContext(ctx, "tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to dial with broker %s: %s", addr, err.Error())
		}
		_, err = conn.ReadPartitions()
		if err != nil {
			return fmt.Errorf("failed to read partitions of broker: %s", err)
		}
	}
	return nil
}
