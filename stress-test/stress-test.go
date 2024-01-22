package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

func main() {
	// Number of concurrent requests
	numRequests := 5

	// Wait group to wait for all Goroutines to finish
	var wg sync.WaitGroup
	wg.Add(numRequests)

	// Launch multiple Goroutines to simulate concurrent requests
	for i := 0; i < numRequests; i++ {
		go func(index int) {
			defer wg.Done()

			// URL of the endpoint you want to hit
			url := "http://localhost:9080"

			// Send a GET request to the endpoint
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Error sending request %d: %v\n", index, err)
				return
			}
			defer resp.Body.Close()

			b, err := io.ReadAll(resp.Body)
			// b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
			if err != nil {
				log.Fatalln(err)
			}

			// Print the response status
			fmt.Printf("Request %d Status: %s Payload: %s\n", index, resp.Status, string(b))
		}(i)
	}

	// Wait for all Goroutines to finish
	wg.Wait()
}
