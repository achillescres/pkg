package mysync

import "sync"

type Once struct {
	once *sync.Once
	done bool
}

func NewOnce() *Once {
	return &Once{once: &sync.Once{}, done: false}
}

func (once *Once) Do(f func()) {
	once.once.Do(func() {
		once.done = true
		f()
	})
}

func (once *Once) IsDone() bool {
	return once.done
}
