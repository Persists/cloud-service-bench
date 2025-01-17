package main

import (
	"cloud-service-bench/internal/archive"
	"cloud-service-bench/internal/config"
	"cloud-service-bench/internal/monitoring"
	"cloud-service-bench/internal/timeline"
	"fmt"
	"os"
	"time"
)

func main() {
	flags, err := config.GetFlags()
	if err != nil {
		fmt.Println(err)
		return
	}

	cfg, err := config.LoadConfig("./config/experiment/config.yml")
	if err != nil {
		fmt.Println(err)
		return
	}

	CPUMonitor := &monitoring.CPUMonitor{}
	MemMonitor := &monitoring.MemMonitor{}
	NetworkMonitor := &monitoring.NetworkMonitor{}

	monitor := monitoring.NewMonitor(CPUMonitor, MemMonitor, NetworkMonitor)

	err = os.MkdirAll(cfg.Archive.Directory, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	metadata := "Job: Monitor\n"
	metadata += config.GenerateMetadata(cfg, flags.InstanceName, flags.Zone)

	fmt.Println(metadata)
	filePath := cfg.Archive.Directory + "/" + fmt.Sprintf("monitor_%s_%s_w%d.log", flags.InstanceName, cfg.Experiment.Id, cfg.Generator.Workers)

	ac, err := archive.NewFileArchiveClient(filePath, metadata)
	if err != nil {
		fmt.Println(err)
		return
	}
	ac.Start()

	t := &timeline.TimeLine{}

	// Set the warm-up, experiment and cool-down phases
	t.SetWarmUp(time.Duration(cfg.Experiment.WarmUp)*time.Second, func() {
		ac.Write("Warm-up phase, starting monitoring at " + time.Now().Format("2006-01-02T15:04:05.000Z") + "\n")
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			stats, err := monitor.GetStats()
			if err != nil {
				fmt.Println(err)
				return
			}

			ac.Write(stats.String())
		}
	})
	t.SetExperiment(time.Duration(cfg.Experiment.Duration)*time.Second, func() {
		ac.Write("Experiment phase, continue monitoring at " + time.Now().Format("2006-01-02T15:04:05.000Z") + "\n")
	})
	t.SetCoolDown(time.Duration(cfg.Experiment.CoolDown)*time.Second, func() {
		ac.Write("Cool-down phase, continue monitoring  at " + time.Now().Format("2006-01-02T15:04:05.000Z") + "\n")
	})

	startAt := flags.StartAt
	if startAt.Before(time.Now()) {
		startAt = time.Now()
	}

	// Run the timeline (warm-up, experiment, cool-down)
	t.Run(startAt)

	ac.Write("Finished at " + time.Now().Format("2006-01-02T15:04:05.000Z") + "\n")
	ac.Close()
}
