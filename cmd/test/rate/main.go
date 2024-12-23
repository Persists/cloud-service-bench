// start tcp server
package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	mu                sync.Mutex
	messageTimestamps = []time.Time{}
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Rate Counter Server started on 8080")

	go RateOrchestrater()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}

}

func RateOrchestrater() {
	ticker := time.NewTicker(1 * time.Second)

	// every second, delete all timestamps older than 5 seconds
	// announce current rate

	for range ticker.C {
		mu.Lock()

		// delete all timestamps older than 5 seconds
		for len(messageTimestamps) > 0 && time.Since(messageTimestamps[0]) > 5*time.Second {
			messageTimestamps = messageTimestamps[1:]
		}

		fmt.Println("Current rate:", len(messageTimestamps)/5.0, "messages per second")

		mu.Unlock()
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		now := time.Now()

		mu.Lock()

		messageTimestamps = append(messageTimestamps, now)

		mu.Unlock()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from connection:", err)
	}
}
