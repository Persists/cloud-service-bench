package generator

import (
	"cloud-service-bench/internal/archive"
	"cloud-service-bench/internal/config"
	"cloud-service-bench/internal/log"
	"fmt"
	"sync"
	"time"
)

type GeneratorClient struct {
	archiver        archive.Archiver
	FluentdConfig   *config.TcpConfig
	LogSynthesizer  *log.LogSynthesizer
	GeneratorConfig *config.GeneratorConfig
}

func NewClient(generatorConfig *config.GeneratorConfig, tcpConfig *config.TcpConfig, ac archive.Archiver) *GeneratorClient {
	client := &GeneratorClient{
		archiver:        ac,
		FluentdConfig:   tcpConfig,
		LogSynthesizer:  log.NewLogSynthesizer(generatorConfig.Name, generatorConfig.MessageLength),
		GeneratorConfig: generatorConfig,
	}

	return client
}

func (g *GeneratorClient) Start() {
	syntheticLogs := g.LogSynthesizer.SynthesizeLogs(g.GeneratorConfig.SampleLength)

	ready := sync.WaitGroup{}
	stop := make(chan struct{})

	ready.Add(g.GeneratorConfig.Workers)
	for i := 0; i < g.GeneratorConfig.Workers; i++ {
		go func() {
			err := routine(
				g.FluentdConfig,
				syntheticLogs,
				float64(g.GeneratorConfig.BatchesPerSec),
				g.GeneratorConfig.LogsPerSecond/g.GeneratorConfig.BatchesPerSec,
				&ready,
				stop,
				g.archiver,
			)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	ready.Wait()

	<-time.After(time.Duration(g.GeneratorConfig.Duration) * time.Second)
	close(stop)

	// TODO: print all sent messages to a file (include metadata (zone, instance, worker, etc))

	// TODO: send a last log, so the http server knows when experiment is finished
	fmt.Println("Finished")
}
