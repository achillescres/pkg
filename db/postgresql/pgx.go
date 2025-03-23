package postgresql

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"os"

	"github.com/achillescres/pkg/utils"
	"github.com/jackc/pgx/v5/stdlib"
)

func NewStdlibDB(ctx context.Context, cc ClientConfig) (*sql.DB, error) {
	ew := utils.NewErrorWrapper("NewStdlibDB")

	config, err := NewConfigFromClientConfig(&cc)
	if err != nil {
		return nil, ew(fmt.Errorf("create config: %w", err))
	}

	if cc.UseCA {
		rootCertPool := x509.NewCertPool()
		pem, err := os.ReadFile(cc.CaAbsPath)
		if err != nil {
			return nil, ew(fmt.Errorf("read ca file: %w", err))
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, ew(fmt.Errorf("append certs to cert pool: %w", err))
		}
		config.ConnConfig.TLSConfig = &tls.Config{
			RootCAs:            rootCertPool,
			InsecureSkipVerify: true,
		}
	}

	db := stdlib.OpenDB(*config.ConnConfig)
	if err != nil {
		return nil, ew(fmt.Errorf("open db: %w", err))
	}
	if db == nil {
		return nil, ew(fmt.Errorf("error couldn't connect to db"))
	}

	pingCtx, cancel := context.WithTimeout(ctx, cc.WaitingDuration)
	defer cancel()

	err = db.PingContext(pingCtx)
	if err != nil {
		return nil, ew(fmt.Errorf("ping db: %w", err))
	}

	return db, nil
}
