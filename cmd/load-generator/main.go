package main

import (
	"cloud-service-bench/internal/config"
	"cloud-service-bench/internal/generator"
	"fmt"
	"os"
	"time"
)

func main() {
	startAtStr := os.Args[1]
	startAt, err := time.Parse(time.RFC3339, startAtStr)
	if err != nil {
		fmt.Println("Invalid start time format:", err)
		return
	}

	config, error := config.LoadConfig("./config/experiment/config.yml")
	if error != nil {
		fmt.Println(error)
		return
	}

	instanceName := os.Getenv("INSTANCE_NAME")
	if instanceName == "" {
		fmt.Println("INSTANCE_NAME environment variable is not set")
		return
	}

	client := generator.NewClient(&config.Generator, &config.Fluentd, instanceName)

	start := time.Now()

	client.Start(startAt)

	elapsed := time.Since(start)

	fmt.Println("Message sent to Fluentd successfully!")
	fmt.Printf("Time taken: %v\n", elapsed)
}
