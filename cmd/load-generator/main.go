package main

import (
	"cloud-service-bench/internal/config"
	"cloud-service-bench/internal/generator"
	"fmt"
	"os"
	"time"
)

func main() {
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

	client.Start()

	elapsed := time.Since(start)

	fmt.Println("Message sent to Fluentd successfully!")
	fmt.Printf("Time taken: %v\n", elapsed)
}
