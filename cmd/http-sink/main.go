package main

import (
	"cloud-service-bench/internal/archive"
	"cloud-service-bench/internal/config"
	"cloud-service-bench/internal/sink"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	instanceName := flag.String("instance-name", "", "The name of the instance")
	zone := flag.String("zone", "europe-west3-c", "The zone of the instance")
	flag.Parse()

	if *instanceName == "" {
		fmt.Println("instance-name flag is not set")
		return
	}

	cfg, error := config.LoadConfig("./config/experiment/config.yml")
	if error != nil {
		fmt.Println(error)
		return
	}

	directory := cfg.Archive.Directory
	fmt.Println("Directory: ", directory)
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	metadata := "Job: Sink\n"
	metadata += config.GenerateMetadata(cfg, *instanceName, *zone)
	filePath := directory + "/" + fmt.Sprintf("%s_%s_%dw_%dlps.log", *instanceName, cfg.Experiment.Id, cfg.Generator.Workers, cfg.Generator.LogsPerSecond)

	ac, err := archive.NewFileArchiveClient(filePath, metadata)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ac.Close()
	ac.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sink := sink.NewHttpSink(ac)

	// effectively an infinite timer
	timer := time.NewTimer(time.Hour * 24 * 365 * 100)

	http.HandleFunc("/fluentd", func(w http.ResponseWriter, r *http.Request) {
		timer.Reset(time.Second * 10)
		sink.Handler(w, r)
	})

	fmt.Println("Server is listening on port", cfg.Sink.Port)
	server := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Sink.Port)}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	select {
	case <-sigChan:
		ac.Write(fmt.Sprintf("Finished at %s, because of a signal", time.Now().Format("2006-01-02T15:04:05.000Z")))
	case <-timer.C:
		ac.Write(fmt.Sprintf("Finished at %s, because the timer expired", time.Now().Format("2006-01-02T15:04:05.000Z")))
	}
	ac.Flush()
}
