package generator

import (
	"cloud-service-bench/internal/log"
)

type GeneratorClient struct {
	LogSynthesizer *log.LogSynthesizer
	SampleLength   int
	LogsPerSecond  int
	Duration       int
}

func NewClient(generatorConfig *GeneratorConfig) *GeneratorClient {
	client := &GeneratorClient{
		LogSynthesizer: log.NewLogSynthesizer(generatorConfig.Name, generatorConfig.MessageLength),
		SampleLength:   generatorConfig.SampleLength,
		LogsPerSecond:  generatorConfig.LogsPerSecond,
		Duration:       generatorConfig.Duration,
	}

	return client
}

func (g *GeneratorClient) Start() {}
