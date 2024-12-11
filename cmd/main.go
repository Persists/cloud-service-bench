package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"math/rand"
)

func main() {
	// Define Fluentd server address and port
	fluentdAddress := "victoria-r1.c.cloud-service-be.internal:20001" // Internal DNS hostname for Fluentd

	// Create a TCP connection to the Fluentd server
	conn, err := net.Dial("tcp", fluentdAddress)
	if err != nil {
		fmt.Printf("Error connecting to Fluentd: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	randomString := func(length int) string {
		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		result := make([]byte, length)
		for i := range result {
			result[i] = charset[rand.Intn(len(charset))]
		}
		return string(result)
	}

	// Message to send to Fluentd
	tag := "test.log"                                           // Replace with your desired tag
	message := fmt.Sprintf(`message: "%s"`, randomString(1000)) // Replace 20 with your desired string length
	// Format the message as Fluentd expects: "<tag>\t<JSON>\n"
	fluentdMessage := fmt.Sprintf("%s\t%s\n", tag, message)

	for j := 0; j < (10 * 1000); j++ {
		// Send the message
		for i := 0; i < 10; i++ {
			_, err = conn.Write([]byte(fluentdMessage))
			if err != nil {
				fmt.Printf("Error sending message to Fluentd: %v\n", err)
				os.Exit(1)
			}
		}
		time.Sleep(1 * time.Millisecond)
	}

	fmt.Println("Message sent to Fluentd successfully!")
}
