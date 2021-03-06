package wait_test

import (
	"fmt"
	"github.com/elgohr/stop-and-go/constraint"
	"github.com/elgohr/stop-and-go/wait"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

func ExampleFor() {
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
		if err != nil {
			log.Fatalln(err)
		}
		w1.Done()
	}()

	fmt.Println(wait.For(
		constraint.NoOrder(w3),
		constraint.Before(w1, w2),
	))
	// Output: <nil>
}

func ExampleFailing() {
	w1 := wait.NewWaiter(time.Second)
	fmt.Println(wait.For(constraint.NoOrder(w1)))
	// Output: failed to wait on waiter 1 of 1
}
