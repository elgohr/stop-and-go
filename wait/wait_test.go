package wait_test

import (
	"github.com/elgohr/stop-and-go/constraint"
	"github.com/elgohr/stop-and-go/wait"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	w1 := wait.NewWaiter(time.Second)
	w2 := wait.NewWaiter(2 * time.Second)
	w3 := wait.NewWaiter(3 * time.Second)
	w4 := wait.NewWaiter(4 * time.Second)
	var (
		calledW1 bool
		calledW2 bool
		calledW3 bool
		calledW4 bool
	)

	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		calledW2 = true
		w2.Done()
	}))
	defer ts1.Close()

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		calledW4 = true
		w4.Done()
	}))
	defer ts2.Close()

	go func() {
		calledW1 = true
		w1.Done()
	}()
	go func() {
		calledW3 = true
		assert.False(t, calledW2)
		_, err := http.Get(ts1.URL)
		require.NoError(t, err)
		_, err = http.Get(ts2.URL)
		require.NoError(t, err)
		w3.Done()
	}()

	require.NoError(t, wait.For(
		constraint.NoOrder(w1),
		constraint.Before(w3, w2),
		constraint.Before(w3, w4),
	))

	require.True(t, calledW1)
	require.True(t, calledW2)
	require.True(t, calledW3)
	require.True(t, calledW4)
}

func TestWait_Errors(t *testing.T) {
	w1 := wait.NewWaiter(time.Millisecond)
	require.EqualError(t, wait.For(constraint.NoOrder(w1)), "failed to wait on waiter 1 of 1")
}

func TestWait_ErrorsWithMultiple(t *testing.T) {
	w1 := wait.NewWaiter(time.Millisecond)
	w2 := wait.NewWaiter(time.Millisecond)
	go func() {
		w1.Done()
	}()
	require.EqualError(t, wait.For(constraint.Before(w1, w2)), "failed to wait on waiter 2 of 2")
}

func TestWait_ErrorsWithMultipleUncalled(t *testing.T) {
	w1 := wait.NewWaiter(time.Millisecond)
	w2 := wait.NewWaiter(time.Millisecond)
	w3 := wait.NewWaiter(time.Millisecond)
	w4 := wait.NewWaiter(time.Millisecond)
	w5 := wait.NewWaiter(time.Millisecond)
	require.EqualError(t, wait.For(
		constraint.Before(w1, w2),
		constraint.Before(w3, w1),
		constraint.Before(w4, w2),
		constraint.Before(w3, w4),
		constraint.Before(w5, w1),
	), "failed to wait on waiter 1 of 5")
}
