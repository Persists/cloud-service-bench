package log

import (
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

func (log *LogMessage) ToArchivable() string {
	return fmt.Sprintf("%s %s %s [%s %s]", log.Time.Format("2006-01-02T15:04:05.000Z"), log.Name, log.Tags, log.Severity, log.Message)
}
