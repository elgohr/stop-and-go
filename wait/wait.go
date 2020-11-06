package wait

import (
	"fmt"
	"time"
)

type Waiter struct {
	timeout time.Duration
	w       chan struct{}
}

type Option func(w []Waiter) []Waiter

func NewWaiter(timeout time.Duration) Waiter {
	return Waiter{
		timeout: timeout,
		w:       make(chan struct{}, 1),
	}
}

func (w *Waiter) Done() {
	w.w <- struct{}{}
}

func For(opts ...Option) error {
	waiters := []Waiter{}
	for _, opt := range opts {
		waiters = opt(waiters)
	}
	for i, w := range waiters {
		select {
		case <-w.w:

		case <-time.Tick(w.timeout):
			return fmt.Errorf("failed to wait on waiter %d of %d", i+1, len(waiters))
		}
	}
	return nil
}
