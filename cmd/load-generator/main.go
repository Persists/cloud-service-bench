package main

import (
	"cloud-service-bench/internal/config"
	"cloud-service-bench/internal/generator"
	"flag"
	"fmt"
	"time"
)

func main() {
	startAtStr := flag.String("start-at", "", "Time to start the load generator in RFC3339 format")
	instanceName := flag.String("instance-name", "", "The name of the instance")
	flag.Parse()

	if *startAtStr == "" {
		fmt.Println("start-at flag is not set")
		return
	}

	if *instanceName == "" {
		fmt.Println("instance-name flag is not set")
		return
	}

	startAt, err := time.Parse(time.RFC3339, *startAtStr)
	if err != nil {
		fmt.Println("Invalid startAt time format. Please use RFC3339 format.")
		return
	}

	config, error := config.LoadConfig("./config/experiment/config.yml")
	if error != nil {
		fmt.Println(error)
		return
	}

	client := generator.NewClient(&config.Generator, &config.Fluentd, *instanceName)

	start := time.Now()

	client.Start(startAt)

	elapsed := time.Since(start)

	fmt.Println("Message sent to Fluentd successfully!")
	fmt.Printf("Time taken: %v\n", elapsed)
}
