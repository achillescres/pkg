package containerutils

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"testing"
)

func CleanupContainerAfterTest(t *testing.T, ctx context.Context, container testcontainers.Container) {
	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})
}
