package log

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type LogMessage struct {
	Time     time.Time `json:"timestamp"`
	Severity Severity  `json:"severity"`
	Name     string    `json:"name"`
	Message  string    `json:"message"`
	Tags     []string  `json:"tags"`
}

type SlimLogMessage struct {
	Time time.Time `json:"timestamp"`
	Name string    `json:"name"`
	Tags []string  `json:"tags"`
}

type Severity string

const (
	INFO  Severity = "INFO"
	WARN  Severity = "WARN"
	ERROR Severity = "ERROR"
	DEBUG Severity = "DEBUG"
)

// CreateFluentdMessage generates a Fluentd-compatible log message string
func (log *LogMessage) ToFluentdMessage() string {
	tags := strings.Join(log.Tags, ",")
	return fmt.Sprintf("%s %s %s [%s] %s\n", log.Time.Format("2006-01-02T15:04:05.000Z"), log.Severity, log.Name, tags, log.Message)
}

// ToArchivable converts a LogMessage to a JSON string
func (log *LogMessage) ToArchivable() (string, error) {
	jsonBytes, err := json.Marshal(log)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil

}
