package mysync

import "sync"

type WaitGroupWithDone struct {
	wg   *sync.WaitGroup
	done chan struct{}
}

func NewWaitGroupWithDone() *WaitGroupWithDone {
	return &WaitGroupWithDone{wg: new(sync.WaitGroup), done: make(chan struct{})}
}

func (awg *WaitGroupWithDone) WaitChan() <-chan struct{} {
	go func() {
		awg.wg.Wait()
		awg.done <- struct{}{}
	}()

	return awg.done
}

func (awg *WaitGroupWithDone) Add(delta int) {
	awg.wg.Add(delta)
}

func (awg *WaitGroupWithDone) Done() {
	awg.wg.Done()
}

// Wait locks goroutine until wg counter zeros
func (awg *WaitGroupWithDone) Wait() {
	awg.wg.Wait()
}
