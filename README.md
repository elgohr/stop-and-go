# stop-and-go
[![Actions Status](https://github.com/elgohr/stop-and-go/workflows/Test/badge.svg)](https://github.com/elgohr/stop-and-go/actions)
[![codecov](https://codecov.io/gh/elgohr/stop-and-go/branch/master/graph/badge.svg)](https://codecov.io/gh/elgohr/stop-and-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/elgohr/stop-and-go)](https://goreportcard.com/report/github.com/elgohr/stop-and-go)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/elgohr/stop-and-go)](https://pkg.go.dev/github.com/elgohr/stop-and-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Testing helper for concurrency

## Install

```bash
go get -u github.com/elgohr/stop-and-go
```

## Usage

```go
package main

import stopandgo "github.com/elgohr/stop-and-go"


func main() {
	w1 := stopandgo.wait.NewWaiter(time.Second)
	w2 := stopandgo.wait.NewWaiter(time.Second)
	w3 := stopandgo.wait.NewWaiter(time.Second)

	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w2.Done()
	}))
	defer ts1.Close()

	go func() {
		w3.Done()
	}()

	go func() {
		if _, err := http.Get(ts1.URL); err != nil {
			Panic(err)
		}
		w1.Done()
	}()

	if err := stopandgo.wait.For(
		stopandgo.constraint.NoOrder(w3),
		stopandgo.constraint.Before(w1, w2),
	); err != nil {
		Panic(err)
	}
}
```
