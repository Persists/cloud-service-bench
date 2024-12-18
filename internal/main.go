// package main

// import (
// 	"cloud-service-bench/internal/connection"
// 	"cloud-service-bench/internal/generator"
// 	"fmt"
// 	"os"
// 	"sync"
// 	"time"
// )

// func routine(ready, wg *sync.WaitGroup, config *generator.Config) {
// 	defer wg.Done()

// 	connectionClient := connection.NewConnectionClient(config.Fluentd.Host, config.Fluentd.Port)
// 	connectionClient.Connect()
// 	defer connectionClient.Disconnect()

// 	start := time.Now()

// 	syntheticLogs := client.LogSynthesizer.SynthesizeLogs(config.Generator.SampleLength / config.Generator.Workers)

// 	ready.Done()
// 	ready.Wait()

// 	elapsed := time.Since(start)
// 	fmt.Printf("Time taken to generate logs: %v\n", elapsed)

// 	fmt.Println("All workers ready, sending logs to Fluentd")

// 	for i := 0; i < len(syntheticLogs); i++ {
// 		err := connectionClient.SendLog(&syntheticLogs[i])
// 		if err != nil {
// 			fmt.Printf("Error sending message to Fluentd: %v\n", err)
// 			os.Exit(2)
// 		}
// 	}
// }

// func main() {
// 	config, error := generator.LoadConfig("./config/experiment/config.yml")

// 	if error != nil {
// 		fmt.Println(error)
// 		os.Exit(1)
// 	}

// 	var wg sync.WaitGroup
// 	var ready sync.WaitGroup
// 	start := time.Now()

// 	ready.Add(config.Generator.Workers)
// 	for i := 0; i < config.Generator.Workers; i++ {
// 		wg.Add(1)
// 		go routine(&ready, &wg, config)
// 	}

// 	wg.Wait()

// 	elapsed := time.Since(start)

// 	fmt.Println("Message sent to Fluentd successfully!")
// 	fmt.Printf("Time taken: %v\n", elapsed)
// }
