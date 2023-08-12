package utils

import "time"

// RunWithTimeout runs f in goroutine and panics if f won't be executed in to
func RunWithTimeout(f func(), to time.Duration) {
	done := make(chan struct{})

	go func() {
		f()
		done <- struct{}{}
	}()

	go func() {
		time.Sleep(to)
		select {
		case <-done:
			return
		default:
		}

		select {
		case <-done:
			return
		default:
			panic("RunWithTimeout: timeout exceed")
		}
	}()
}
