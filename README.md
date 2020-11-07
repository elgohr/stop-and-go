# stop-and-go
[![Actions Status](https://github.com/elgohr/stop-and-go/workflows/Test/badge.svg)](https://github.com/elgohr/stop-and-go/actions)
[![codecov](https://codecov.io/gh/elgohr/stop-and-go/branch/master/graph/badge.svg)](https://codecov.io/gh/elgohr/stop-and-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/elgohr/stop-and-go)](https://goreportcard.com/report/github.com/elgohr/stop-and-go)
[![](https://godoc.org/github.com/elgohr/stop-and-go?status.svg)](http://godoc.org/github.com/elgohr/stop-and-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Testing helper for concurrency

## Install

```bash
go get -u github.com/elgohr/stop-and-go
```

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