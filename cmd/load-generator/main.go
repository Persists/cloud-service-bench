package main

import (
	"cloud-service-bench/internal/config"
	"cloud-service-bench/internal/generator"
	"cloud-service-bench/internal/timeline"
	"fmt"
	"log"
	"time"
)

func main() {
	flags, err := config.GetFlags()
	if err != nil {
		fmt.Println(err)
		return
	}

	cfg, error := config.LoadConfig("./config/experiment/config.yml")
	if error != nil {
		fmt.Println(error)
		return
	}

	metadata := "Job: Generator\n"
	metadata += config.GenerateMetadata(cfg, flags.InstanceName, flags.Zone)

	fmt.Println(metadata)

	// Initialize the generator client
	client := generator.NewClient(&cfg.Generator, &cfg.Experiment, &cfg.Fluentd, flags.InstanceName)

	t := &timeline.TimeLine{}

	// Set the warm-up, experiment and cool-down phases
	t.SetWarmUp(time.Duration(cfg.Experiment.WarmUp)*time.Second, func() {
		client.Start()
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			log.Printf("Total messages sent: %d\n", client.GetTotalMessages())
		}
	})
	t.SetExperiment(time.Duration(cfg.Experiment.Duration)*time.Second, func() {})
	t.SetCoolDown(time.Duration(cfg.Experiment.CoolDown)*time.Second, client.Stop)

	startAt := flags.StartAt
	if startAt.Before(time.Now()) {
		startAt = time.Now()
	}

	// Run the timeline (warm-up, experiment, cool-down)
	t.Run(startAt)
}
