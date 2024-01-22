package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hadiweslati/simplesurance-go-challenge/counter"
)

var semaphore = make(chan struct{}, 1)

func main() {
	// Load previous state from file
	c := counter.NewCountRequest()

	// Start a goroutine to reset the counter every 60 seconds
	go func() {
		for {
			time.Sleep(time.Second * 60)
			c.ResetCounter()
			c.SaveState()
		}
	}()

	// Setup HTTP handlers with middleware
	http.HandleFunc("/", middleware(requestHandler))

	// Start the HTTP server
	fmt.Println("Server listening on :9080")
	err := http.ListenAndServe(":9080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Acquire the semaphore, allowing only one concurrent request
		semaphore <- struct{}{}
		defer func() {
			// Release the semaphore when the function exits
			<-semaphore
		}()

		// Load previous state from file
		requestCounter := counter.NewCountRequest()

		requestCounter.LoadState()
		requestCounter.Count++

		// Save State to avoid using in-memory persist
		requestCounter.SaveState()

		// Pass the request counter to the handler using context
		ctx := context.WithValue(r.Context(), "requestCounter", requestCounter.Count)
		next(w, r.WithContext(ctx))
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the request counter from the context
	ctx := r.Context()
	count, ok := ctx.Value("requestCounter").(int)
	if !ok {
		// Handle the case where the value is not found in the context
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Display the current request count
	fmt.Fprintf(w, "Total Requests in Last 60 Seconds: %d", count)
}
