package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// Define Fluentd server address and port
	fluentdAddress := "34.40.30.116:20001" // Replace with your Fluentd host and port

	// Create a TCP connection to the Fluentd server
	conn, err := net.Dial("tcp", fluentdAddress)
	if err != nil {
		fmt.Printf("Error connecting to Fluentd: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Message to send to Fluentd
	tag := "test.log" // Replace with your desired tag
	message := `message:test`

	// Format the message as Fluentd expects: "<tag>\t<JSON>\n"
	fluentdMessage := fmt.Sprintf("%s\t%s\n", tag, message)

	// Send the message
	_, err = conn.Write([]byte(fluentdMessage))
	if err != nil {
		fmt.Printf("Error sending message to Fluentd: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Message sent to Fluentd successfully!")
}
