# stop-and-go
Testing helper for concurrency

## Usage

```go
import "github.com/stretchr/testify/require"

func TestExample(t *testing.T) {
	w1 := wait.NewWaiter(time.Second)
	w2 := wait.NewWaiter(time.Second)
	w3 := wait.NewWaiter(time.Second)

	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w2.Done()
	}))
	defer ts1.Close()

	go func() {
		w3.Done()
	}()

	go func() {
		_, err := http.Get(ts1.URL)
		require.NoError(t, err)
		w1.Done()
	}()

	require.NoError(t, wait.For(
		constraint.NoOrder(w3),
		constraint.Before(w1, w2),
	))
}
```