package main

import (
	"fmt"
	"math"
	"net"
	"time"
)

const (
	// number of messages per worker
	numOfMessages = 1430
	workers       = 1
)

func main() {
	var tickerInterval float64
	var messagesPerTick int
	highestDevisor := findHighestDivisor(numOfMessages)
	if highestDevisor <= 9 {
		tickerInterval, messagesPerTick = findEfficientRate(numOfMessages, 10.0)
	} else {
		tickerInterval = float64(highestDevisor)
		messagesPerTick = numOfMessages / highestDevisor
	}

	println("Ticker interval:", tickerInterval)
	println("Messages per tick:", messagesPerTick)

	for i := 0; i < workers; i++ {
		go worker(tickerInterval, messagesPerTick)
	}

	select {}
}

func worker(tickerInterval float64, messagesPerTick int) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(time.Second / time.Duration(tickerInterval))
	defer ticker.Stop()

	for range ticker.C {
		for j := 0; j < messagesPerTick; j++ {
			message := fmt.Sprintf("%s INFO test [tag1,tag2] This is a test message adsa sd ad ad ad adsd asd ", time.Now().Format(time.RFC3339))
			_, err := fmt.Fprintln(conn, message)
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
		}
	}
}

func findHighestDivisor(rate int) int {
	highestDivisor := 0

	// Iterate from 100 down to 1
	for d := 100; d > 1; d-- {
		if rate%d == 0 {
			highestDivisor = d
			break
		}
	}

	if highestDivisor == 0 {
		return 0
	}

	return highestDivisor
}

// findEfficientRate finds the most efficient rate (how many ticks per second and how many messages per tick)
// for a given number of messages per second. Therefore it approximates the rate to the closest
func findEfficientRate(rate int, limit float64) (float64, int) {
	tickerInterval := 0.0
	messagesPerTick := 1
	minDiff := math.MaxFloat64

	for a := 1; a <= rate; a++ {
		y := float64(rate) / float64(a)
		diff := math.Abs(y - limit)

		if diff < minDiff {
			minDiff = diff
			tickerInterval = y
			messagesPerTick = a
		}

		if y < limit {
			break
		}
	}

	return tickerInterval, messagesPerTick
}
