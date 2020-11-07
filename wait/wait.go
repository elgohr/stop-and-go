package wait

import (
	"fmt"
	"time"
)

// Waiter represents a point in tests to wait for
type Waiter struct {
	timeout time.Duration
	w       chan struct{}
}

// NewWaiter constructs a new Waiter
// Needs a timeout, which is the longest time to wait for the Waiter
func NewWaiter(timeout time.Duration) Waiter {
	return Waiter{
		timeout: timeout,
		w:       make(chan struct{}, 1),
	}
}

// Done marks the Waiter as called
func (w *Waiter) Done() {
	w.w <- struct{}{}
}

// Option is used to configure the dependencies for Waiter
type Option func(w []Waiter) []Waiter

// For provides a way to configure dependencies between Waiters
// It errors when at least one Waiter hasn't been called
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
