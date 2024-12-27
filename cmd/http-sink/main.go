package main

import (
	"cloud-service-bench/internal/archive"
	"cloud-service-bench/internal/config"
	"cloud-service-bench/internal/sink"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, error := config.LoadConfig("./config/experiment/config.yml")
	if error != nil {
		fmt.Println(error)
		return
	}

	instanceName := config.GetEnv("INSTANCE_NAME")
	zone := config.GetEnv("ZONE")

	ac, err := archive.NewFileArchiveClient(cfg, instanceName, zone)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ac.Close()
	ac.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("Received interrupt signal, flushing archive client...")
		ac.Flush()
		os.Exit(0)
	}()

	sink := sink.NewHttpSink(ac)

	http.HandleFunc("/fluentd", sink.Handler)

	fmt.Println("Server is listening on port", cfg.Sink.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Sink.Port), nil))
}
