package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"os"
	"time"
)

type DialConfig struct {
	Username, Password string
	UseSASL            bool
	UseCA              bool
	CaAbsPath          string
	Timeout            time.Duration
}

func NewDialer(config DialConfig) (*kafka.Dialer, error) {
	dialer := &kafka.Dialer{
		Timeout:   config.Timeout,
		DualStack: true,
	}

	if config.UseSASL {
		mechanism, err := scram.Mechanism(scram.SHA512, config.Username, config.Password)
		if err != nil {
			return nil, fmt.Errorf("create mechanism: %w", err)
		}
		dialer.SASLMechanism = mechanism
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

		dialer.TLS = &tls.Config{
			RootCAs:            rootCertPool,
			InsecureSkipVerify: true,
		}
	}

	return dialer, nil
}
