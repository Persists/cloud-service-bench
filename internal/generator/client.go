package generator

import (
	"cloud-service-bench/internal/config"
	"cloud-service-bench/internal/log"
	"fmt"
)

type GeneratorClient struct {
	ExperimentConfig *config.ExperimentConfig
	workers          []*Worker
	startChan        chan struct{}
	stopChan         chan struct{}
}

func (g *GeneratorClient) GetTotalMessages() int {
	totalMessages := 0
	for _, worker := range g.workers {
		totalMessages += worker.TotalMessages
	}
	return totalMessages
}

// NewClient creates a new instance of GeneratorClient with the provided configurations.
// It initializes the workers, log synthesizer, and communication channels.
//
// It starts the workers by newWorker.
//
// Parameters:
//   - generatorConfig: Configuration for the log generator, including the number of workers, message length, and sample length.
//   - experimentConfig: Configuration for the experiment.
//   - tcpConfig: Configuration for TCP connections.
//   - name: Name identifier for the log synthesizer.
func NewClient(generatorConfig *config.GeneratorConfig, experimentConfig *config.ExperimentConfig, tcpConfig *config.TcpConfig, name string) *GeneratorClient {
	workers := make([]*Worker, generatorConfig.Workers)
	logSynthesizer := log.NewLogSynthesizer(name, generatorConfig.MessageLength)
	syntheticLogs := logSynthesizer.SynthesizeLogs(generatorConfig.SampleLength)
	startChan := make(chan struct{})
	stopChan := make(chan struct{})

	for i := 0; i < generatorConfig.Workers; i++ {
		worker, err := newWorker(tcpConfig, syntheticLogs, startChan, stopChan, i)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		workers[i] = worker
	}

	return &GeneratorClient{
		ExperimentConfig: experimentConfig,
		workers:          workers,
		startChan:        startChan,
		stopChan:         stopChan,
	}
}

// Start starts the workers
func (g *GeneratorClient) Start() {
	close(g.startChan)
}

// Stop stops the workers
func (g *GeneratorClient) Stop() {
	close(g.stopChan)
}
