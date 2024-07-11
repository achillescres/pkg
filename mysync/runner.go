package mysync

import (
	"context"
	"sync"
)

// Runner runs only one function in one time
type Runner struct {
	mu      sync.Mutex
	errChan chan error
}

func NewRunner() *Runner {
	return &Runner{mu: sync.Mutex{}, errChan: make(chan error, 1)}
}

// RunWithRecovery locks current goroutine then execute given function in new goroutine
// and when it returns error unlocks current goroutine and returns error
func (s *Runner) Run(f func() error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return f()
}

// RunWithRecovery runs f on error calls callback on panic recovers and continues to run f
func RunWithRecovery(ctx context.Context, f func(ctx context.Context) error, recoverCallback func(error)) error {
	err := make(chan error, 1)
	recov := make(chan error, 1)
	for {
		GoardWithCallback(
			ctx,
			func(ctx context.Context) {
				err <- f(ctx)
			},
			func(err error) {
				recov <- err
			},
		)
		select {
		case e := <-err:
			return e
		case err := <-recov:
			if err != nil {
				recoverCallback(err)
			}
		}
	}
}
