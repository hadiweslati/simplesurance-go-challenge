package counter

import (
	"encoding/gob"
	"fmt"
	"os"
	"sync"
	"time"
)

const dataFile = "request_counter.gob"

// Counter struct to store request count and timestamp
type Counter struct {
	Count     int
	ResetTime time.Time
	mu        sync.Mutex
}

func (c *Counter) ResetCounter() {
	// Lock to ensure atomic access to the counter
	c.mu.Lock()
	defer c.mu.Unlock()
	// Reset the counter
	c.Count = 0
	c.ResetTime = time.Now()
}

func (c *Counter) LoadState() {
	// Attempt to open the data file
	file, err := os.Open(dataFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Decode the Counter struct from the file
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		fmt.Println("Error decoding data:", err)
	}
}

func (c *Counter) SaveState() {
	// Attempt to create or open the data file
	file, err := os.Create(dataFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Encode the Counter struct and write to the file
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(&c)
	if err != nil {
		fmt.Println("Error encoding data:", err)
	}
}

func (c *Counter) GetRequestCount() int {
	// Lock to ensure atomic access to the counter
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Count
}

func NewCountRequest() *Counter {
	return &Counter{}
}
