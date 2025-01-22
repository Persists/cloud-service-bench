package main

import (
	"cloud-service-bench/internal/archive"
	"cloud-service-bench/internal/config"
	"cloud-service-bench/internal/sink"
	"cloud-service-bench/internal/timeline"
	"fmt"
	"log"
	"net/http"
	"os"
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

	directory := cfg.Archive.Directory
	fmt.Println("Directory: ", directory)
	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	metadata := "Job: Sink\n"
	metadata += config.GenerateMetadata(cfg, flags.InstanceName, flags.Zone)

	fmt.Println(metadata)
	filePath := directory + "/" + fmt.Sprintf("%s_%s_%dw.log", flags.InstanceName, cfg.Experiment.Id, cfg.Generator.Workers)

	ac, err := archive.NewFileArchiveClient(filePath, metadata)
	if err != nil {
		fmt.Println(err)
		return
	}
	ac.Start()

	sink := sink.NewHttpSink(ac)

	http.HandleFunc("/fluentd", sink.Handler)

	fmt.Printf("Starting sink on port %d\n", cfg.Sink.Port)
	server := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Sink.Port)}

	t := &timeline.TimeLine{}

	// Set the warm-up, experiment and cool-down phases
	t.SetWarmUp(time.Duration(cfg.Experiment.WarmUp)*time.Second, func() {
		ac.Write("Warm-up phase, starting sink at " + time.Now().Format("2006-01-02T15:04:05.000Z") + "\n")
		log.Fatal(server.ListenAndServe())
	})
	t.SetExperiment(time.Duration(cfg.Experiment.Duration)*time.Second, func() {
		ac.Write("\n" + "Experiment phase, continue sink at " + time.Now().Format("2006-01-02T15:04:05.000Z") + "\n")
	})
	t.SetCoolDown(time.Duration(cfg.Experiment.CoolDown)*time.Second, func() {
		ac.Write("\n" + "Cool-down phase, continue sink at " + time.Now().Format("2006-01-02T15:04:05.000Z") + "\n")
	})

	startAt := flags.StartAt
	if startAt.Before(time.Now()) {
		startAt = time.Now()
	}

	// Run the timeline (warm-up, experiment, cool-down)
	t.Run(startAt)

	ac.Close()
}
