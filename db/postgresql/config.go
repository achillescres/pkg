package postgresql

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type ClientConfig struct {
	MaxConnections                               int
	MaxConnectionAttempts                        int
	WaitingDuration                              time.Duration
	Username, Password, Host, Port, DatabaseName string
	UseCA                                        bool
	CaAbsPath                                    string
	SimpleQueryMode                              bool
}

func NewConfigFromClientConfig(cc *ClientConfig) (*pgxpool.Config, error) {
	connstring := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s target_session_attrs=read-write ",
		cc.Host, cc.Port, cc.DatabaseName, cc.Username, cc.Password,
	)
	if cc.UseCA {
		connstring += "sslmode=verify-full"
	}

	config, err := pgxpool.ParseConfig(connstring)
	if err != nil {
		return nil, err
	}

	addOptionsToConfig(cc, config, cc.SimpleQueryMode)

	return config, nil
}

func addOptionsToConfig(cc *ClientConfig, config *pgxpool.Config, simpleProtocol bool) *pgxpool.Config {
	config.MaxConns = int32(cc.MaxConnections)
	if simpleProtocol {
		config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	}
	return config
}
