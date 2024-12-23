package main

import (
	"cloud-service-bench/internal/archive"
	"cloud-service-bench/internal/config"
	"cloud-service-bench/internal/generator"
	"fmt"
	"time"
)

func main() {
	config, error := config.LoadConfig("./config/experiment/config.yml")
	if error != nil {
		fmt.Println(error)
		return
	}

	ac := archive.NewFileArchiveClient(config, "./results", true)
	go ac.Start()
	client := generator.NewClient(&config.Generator, &config.Fluentd, ac)

	start := time.Now()

	client.Start()

	elapsed := time.Since(start)

	fmt.Println("Message sent to Fluentd successfully!")
	fmt.Printf("Time taken: %v\n", elapsed)
}
