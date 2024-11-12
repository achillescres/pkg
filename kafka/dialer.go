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

const ( // shit but reverse compatible shit
	SHA256 uint8 = iota
	SHA512
)

type DialerConfig struct {
	Username, Password      string
	UseSASL                 bool
	UseCA                   bool
	CaAbsPath               string
	Timeout                 time.Duration
	VerifyBrokerCertificate bool
	// Default 0 is SHA256
	ScramAlgorithm uint8
}

func (dc DialerConfig) validate() error {
	ew := utils.NewErrorWrapper("DialerConfig - validate")

	if dc.ScramAlgorithm > 1 {
		return ew(fmt.Errorf("invalid ScramAlgorithm, valid range [0, 1]: %d", dc.ScramAlgorithm))
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
		algo := scram.SHA256
		if config.ScramAlgorithm == 1 {
			algo = scram.SHA512
		}
		mechanism, err := scram.Mechanism(algo, config.Username, config.Password)
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
