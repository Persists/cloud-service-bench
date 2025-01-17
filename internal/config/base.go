package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ExperimentConfig struct {
	Id       string `yaml:"id"`
	Duration int    `yaml:"duration"`
	WarmUp   int    `yaml:"warmup"`
	CoolDown int    `yaml:"cooldown"`
}

type TcpConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type GeneratorConfig struct {
	MessageLength int `yaml:"message-length"`
	SampleLength  int `yaml:"sample-length"`
	Workers       int `yaml:"workers"`
}

type ArchiveConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Directory string `yaml:"directory"`
}

type Sink struct {
	Port int `yaml:"port"`
}

type Config struct {
	Experiment ExperimentConfig `yaml:"experiment"`
	Fluentd    TcpConfig        `yaml:"fluentd"`
	Generator  GeneratorConfig  `yaml:"generator"`
	Archive    ArchiveConfig    `yaml:"archive"`
	Sink       Sink             `yaml:"sink"`
}

func validtateConfig(cfg *Config) error {

	if cfg.Experiment.Id == "" {
		return fmt.Errorf("experiment.id must be provided")
	}

	if cfg.Experiment.Duration <= 0 {
		return fmt.Errorf("experiment.duration must be greater than 0")
	}

	if cfg.Experiment.WarmUp < 0 {
		return fmt.Errorf("experiment.warmup must be greater than or equal to 0")
	}

	if cfg.Experiment.CoolDown < 0 {
		return fmt.Errorf("experiment.cooldown must be greater than or equal to 0")
	}

	if cfg.Generator.MessageLength <= 0 {
		return fmt.Errorf("message-length must be greater than 0")
	}

	if cfg.Generator.SampleLength <= 0 {
		return fmt.Errorf("sample-length must be greater than 0")
	}

	if cfg.Generator.Workers <= 0 {
		return fmt.Errorf("workers must be greater than 0")
	}

	if cfg.Fluentd.Host == "" {
		return fmt.Errorf("fluentd.host must be provided")
	}

	if cfg.Fluentd.Port <= 0 {
		return fmt.Errorf("fluentd.port must be specified")
	}

	if cfg.Archive.Enabled {
		if cfg.Archive.Directory == "" {
			return fmt.Errorf("archive.directory must be provided")
		}
	}

	if cfg.Sink.Port <= 0 {
		return fmt.Errorf("sink.port must be specified")
	}

	return nil
}

// LoadSenderConfig reads a YAML file and unmarshals it into a Config struct
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	if err := validtateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func GenerateMetadata(config *Config, instanceName, zone string) string {
	return fmt.Sprintf(
		"Experiment ID: %s\nInstance Name: %s\nZone: %s\nWorkers: %d\nDuration: %d seconds\nMessage Length: %d\nSample Length: %d\n",
		config.Experiment.Id,
		instanceName,
		zone,
		config.Generator.Workers,
		config.Experiment.Duration,
		config.Generator.MessageLength,
		config.Generator.SampleLength,
	)

}
