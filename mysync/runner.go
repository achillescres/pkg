package mysync

import (
	"sync"
)

// Runner runs only one function in one time
type Runner struct {
	mu      sync.Mutex
	errChan chan error
}

func NewRunner() *Runner {
	return &Runner{mu: sync.Mutex{}, errChan: make(chan error)}
}

// Run locks current goroutine then execute given function in new goroutine
// and when it returns error unlocks current goroutine and returns error
func (s *Runner) Run(f func() error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.errChan <- f()

	return <-s.errChan
}
