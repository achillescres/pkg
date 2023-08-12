package credVault

import (
	"context"
	"time"
)

type PullCred[Cred any] func(ctx context.Context) (Cred, *time.Time, error)
