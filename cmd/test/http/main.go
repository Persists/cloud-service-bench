package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	fluentdURL    = "http://localhost:8080/http.tag"
	totalRequests = 10000
	concurrency   = 100 // Number of concurrent workers
)

func main() {
	var wg sync.WaitGroup
	requests := make(chan int, totalRequests)

	// Start worker goroutines
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range requests {
				if err := sendHTTPRequest(); err != nil {
					fmt.Println("Error sending request:", err)
				}
			}
		}()
	}

	// Enqueue requests
	for i := 0; i < totalRequests; i++ {
		requests <- i
	}
	close(requests)

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("Completed sending all requests.")
}

func sendHTTPRequest() error {
	// Create the payload
	payload := `{"foo":"bar"}`

	// Set the custom timestamp
	timestamp := time.Now().UnixNano() / int64(time.Second)
	timestampStr := strconv.FormatInt(timestamp, 10) + "." + strconv.FormatInt(time.Now().UnixNano()%int64(time.Second), 10)

	// Construct the URL with the timestamp as a query parameter
	u, err := url.Parse(fluentdURL)
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("time", timestampStr)
	u.RawQuery = q.Encode()

	// Send the HTTP POST request
	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
