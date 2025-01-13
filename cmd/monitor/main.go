package main

import (
	"cloud-service-bench/internal/archive"
	"cloud-service-bench/internal/config"
	"cloud-service-bench/internal/monitoring"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig("./config/experiment/config.yml")
	if err != nil {
		fmt.Println(err)
		return
	}

	instanceName := config.GetEnv("INSTANCE_NAME")
	zone := config.GetEnv("ZONE")

	CPUMonitor := &monitoring.CPUMonitor{}
	MemMonitor := &monitoring.MemMonitor{}
	NetworkMonitor := &monitoring.NetworkMonitor{}

	monitor := monitoring.NewMonitor(CPUMonitor, MemMonitor, NetworkMonitor)

	ticker := time.NewTicker(1 * time.Second)

	directory := "./results"
	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	metadata := "Job: Monitor\n"
	metadata += config.GenerateMetadata(cfg, instanceName, zone)
	filePath := directory + "/" + fmt.Sprintf("monitor_%s_%s_%dlps.log", instanceName, cfg.Experiment.Id, cfg.Generator.LogsPerSecond)

	ac, err := archive.NewFileArchiveClient(filePath, metadata)
	if err != nil {
		fmt.Println(err)
		return
	}
	ac.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("Received interrupt signal, flushing archive client...")
		ac.Flush()
		os.Exit(0)
	}()

	for range ticker.C {
		stats, err := monitor.GetStats()
		if err != nil {
			fmt.Println(err)
			return
		}

		ac.Write(stats.String())
		ac.Flush()
	}
}
