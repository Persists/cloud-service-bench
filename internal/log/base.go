package log

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type LogMessage struct {
	Time     time.Time
	Severity Severity
	Name     string
	Message  string
	Tags     []string
}

type Severity int

const (
	INFO Severity = iota
	WARN
	ERROR
	DEBUG
)

func (s Severity) String() string {
	return [...]string{"INFO", "WARN", "ERROR", "DEBUG"}[s]
}

var logMessageRegex = regexp.MustCompile(`^(?P<time>[^\]]+) (?P<severity>[^ ]+) (?P<name>[^ ]+) \[(?P<tags>[^\]]+)\] (?P<message>.+)$`)

// CreateFluentdMessage generates a Fluentd-compatible log message string
func (log *LogMessage) ToFluentdMessage() string {
	tags := strings.Join(log.Tags, ",")
	return fmt.Sprintf("%s %s %s [%s] %s\n", log.Time.Format("2006-01-02T15:04:05.000Z"), log.Severity.String(), log.Name, tags, log.Message)
}
