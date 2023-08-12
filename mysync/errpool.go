package mysync

import "sync"

type ErrorPool struct {
	mu   sync.Mutex
	pool []error
}

func (e *ErrorPool) Throw(err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.pool = append(e.pool, err)
}

func (e *ErrorPool) Grab() []error {
	var pool []error
	copy(e.pool, pool)
	return pool
}
