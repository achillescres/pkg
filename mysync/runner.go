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

// Run locks current goroutine then execute given function in new goroutine
// and when it returns error unlocks current goroutine and returns error
func (s *Runner) Run(f func() error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return f()
}

// Run infinitely runs f on error calls callback on panic recovers and continues to run f
func Run(ctx context.Context, f func(ctx context.Context) error, errCallback func(error)) {
	recov := make(chan error, 1)
	for {
		GoardWithCallback(
			ctx,
			func(ctx context.Context) {
				recov <- f(ctx)
			},
			func(err error) {
				recov <- err
			},
		)
		err := <-recov
		if err != nil {
			errCallback(err)
		}
	}
}
