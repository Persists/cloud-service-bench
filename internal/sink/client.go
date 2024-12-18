package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	// Define the HTTP server port
	http.HandleFunc("/fluentd", func(w http.ResponseWriter, r *http.Request) {
		// Only accept POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusInternalServerError)
			return
		}

		fmt.Println("Request received!")

		// Count the number of log entries in the request body
		go countEntries(body)

		// Respond to the client
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Request received successfully")
	})

	// Start the HTTP server on port 8080
	fmt.Println("Server is listening on port 8080...\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// countEntries estimates the number of entries by counting occurrences of "tag"
func countEntries(body []byte) {
	count := 0
	tagKeyword := `"tag":`
	for i := 0; i < len(body)-len(tagKeyword)+1; i++ {
		if string(body[i:i+len(tagKeyword)]) == tagKeyword {
			count++
		}
	}
	fmt.Printf("Estimated %d entries in this request.\n", count)
}
