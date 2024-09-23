package postgresutils

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"
	"time"
)

type Info struct {
	Image        string
	Host         string
	ExposedPort  nat.Port
	InternalPort nat.Port
	DbName       string
	User         string
	Password     string
}

const (
	DefaultUser     = "test"
	DefaultPassword = "test"
	DefaultDb       = "test"
	DefaultImage    = "postgres:15-alpine"
)

func NewContainer(ctx context.Context, useRyuk bool, opts ...testcontainers.ContainerCustomizer) (*postgres.PostgresContainer, Info, error) {
	if useRyuk {
		err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		if err != nil {
			return nil, Info{}, fmt.Errorf("set TESTCONTAINERS_RYUK_DISABLED to true: %w", err)
		}
	}

	opts = append([]testcontainers.ContainerCustomizer{
		postgres.WithUsername(DefaultUser),
		postgres.WithPassword(DefaultPassword),
		postgres.WithDatabase(DefaultDb),
		postgres.WithInitScripts("./../../../migrations/up.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5 * time.Second)),
	}, opts...)

	requestDoll := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			ExposedPorts: []string{"5432/tcp"},
			Env:          map[string]string{},
			Cmd:          []string{"postgres", "-c", "fsync=off"},
		},
	}
	for _, opt := range opts {
		err := opt.Customize(&requestDoll)
		if err != nil {
			return nil, Info{}, fmt.Errorf("set option on cointainer request: %w", err)
		}
	}

	postgresContainer, err := postgres.Run(ctx, DefaultImage, opts...)
	if err != nil {
		return nil, Info{}, fmt.Errorf("failed to start container: %s", err)
	}

	postgresExposedPort, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, Info{}, fmt.Errorf("get postgres container exposed port: %s", err)
	}

	return postgresContainer,
		Info{
			Image:        requestDoll.Image,
			Host:         "localhost",
			ExposedPort:  postgresExposedPort,
			InternalPort: "5432/tcp",
			DbName:       requestDoll.Env["POSTGRES_DB"],
			User:         requestDoll.Env["POSTGRES_USER"],
			Password:     requestDoll.Env["POSTGRES_PASSWORD"],
		}, nil
}

func RestoreAfterTest(t *testing.T, ctx context.Context, postgresC *postgres.PostgresContainer, snapshotName string) {
	t.Cleanup(func() {
		err := postgresC.Restore(ctx, postgres.WithSnapshotName(snapshotName))
		if err != nil {
			t.Fatalf("restore postgres container to initial snapshot: %s", err)
		}
	})
}
