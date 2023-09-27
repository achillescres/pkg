package postgresql

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"os"
)

type PGXPool struct {
	*pgxpool.Pool
}

func NewPGXPool(ctx context.Context, cc *ClientConfig, logger *logrus.Entry) (PGXPool, error) {
	config, err := NewConfigFromClientConfig(cc)
	if err != nil {
		return PGXPool{}, err
	}

	if cc.UseCA {
		rootCertPool := x509.NewCertPool()
		pem, err := os.ReadFile(cc.CaAbsPath)
		if err != nil {
			return PGXPool{}, err
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return PGXPool{}, fmt.Errorf("error couldn't append certs to cert pool")
		}
		config.ConnConfig.TLSConfig = &tls.Config{
			RootCAs:            rootCertPool,
			InsecureSkipVerify: true,
		}
	}

	logger.Infof("trying to connect to db: %s\n", config.ConnString())
	createCtx, cancel := context.WithTimeout(ctx, cc.WaitingDuration)
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	pool, err := pgxpool.NewWithConfig(createCtx, config)
	if err != nil {
		return PGXPool{}, err
	}
	if pool == nil {
		return PGXPool{}, fmt.Errorf("error couldn't connect to db")
	}

	pingCtx, cancel := context.WithTimeout(ctx, cc.WaitingDuration)
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()
	err = pool.Ping(pingCtx)
	if err != nil {
		return PGXPool{}, err
	}

	return PGXPool{pool}, nil
}
