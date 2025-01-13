package generator

import (
	"cloud-service-bench/internal/config"
	"cloud-service-bench/internal/log"
	"fmt"
	"sync"
	"time"
)

type GeneratorClient struct {
	TCPConfig       *config.TcpConfig
	LogSynthesizer  *log.LogSynthesizer
	GeneratorConfig *config.GeneratorConfig
}

func NewClient(generatorConfig *config.GeneratorConfig, tcpConfig *config.TcpConfig, name string) *GeneratorClient {
	client := &GeneratorClient{
		TCPConfig:       tcpConfig,
		LogSynthesizer:  log.NewLogSynthesizer(name, generatorConfig.MessageLength),
		GeneratorConfig: generatorConfig,
	}

	return client
}

// Start starts the generator client.
// It synthesizes logs and starts the worker routines.
func (g *GeneratorClient) Start(startAt time.Time) {
	syntheticLogs := g.LogSynthesizer.SynthesizeLogs(g.GeneratorConfig.SampleLength)

	ready := sync.WaitGroup{}
	ready.Add(1)
	stop := make(chan struct{})

	ready.Add(g.GeneratorConfig.Workers)
	for i := 0; i < g.GeneratorConfig.Workers; i++ {
		go func() {
			err := g.routine(
				i,
				syntheticLogs,
				&ready,
				stop,
			)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	go func() {
		fmt.Println("Waiting for start time:", time.Until(startAt))
		<-time.After(time.Until(startAt))
		ready.Done()
	}()

	ready.Wait()

	fmt.Println("Started")

	<-time.After(time.Duration(g.GeneratorConfig.Duration) * time.Second)
	close(stop)

	fmt.Println("Finished")
}
