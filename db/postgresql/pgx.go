package postgresql

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/achillescres/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type StdlibConfig struct {
	DSN             string
	UseCA           bool
	CaAbsPath       string
	WaitingDuration time.Duration
}

func NewStdlibDB(ctx context.Context, config StdlibConfig) (*sql.DB, error) {
	ew := utils.NewErrorWrapper("NewStdlibDB")

	pgxConfig, err := pgx.ParseConfig(config.DSN)
	if err != nil {
		return nil, ew(fmt.Errorf("parse config: %w", err))
	}

	if config.UseCA {
		rootCertPool := x509.NewCertPool()
		pem, err := os.ReadFile(config.CaAbsPath)
		if err != nil {
			return nil, ew(fmt.Errorf("read ca file: %w", err))
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, ew(fmt.Errorf("append certs to cert pool: %w", err))
		}
		pgxConfig.TLSConfig = &tls.Config{
			RootCAs:            rootCertPool,
			InsecureSkipVerify: true,
		}
	}

	db := stdlib.OpenDB(*pgxConfig)
	if err != nil {
		return nil, ew(fmt.Errorf("open db: %w", err))
	}
	if db == nil {
		return nil, ew(fmt.Errorf("error couldn't connect to db"))
	}

	pingCtx, cancel := context.WithTimeout(ctx, config.WaitingDuration)
	defer cancel()

	err = db.PingContext(pingCtx)
	if err != nil {
		return nil, ew(fmt.Errorf("ping db: %w", err))
	}

	return db, nil
}
