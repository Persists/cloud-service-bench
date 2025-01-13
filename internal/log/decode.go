package log

import (
	"encoding/json"
	"time"
)

type SlimLogMessage struct {
	Time time.Time `json:"timestamp"`
	Name string    `json:"name"`
	Tags []string  `json:"tags"`
}

func DecodeJson(body []byte) ([]LogMessage, error) {
	var logMessages []LogMessage
	if err := json.Unmarshal(body, &logMessages); err != nil {
		return nil, err
	}
	return logMessages, nil
}

func (log *LogMessage) ToSlimLogMessage() SlimLogMessage {
	return SlimLogMessage{
		Time: log.Time,
		Name: log.Name,
		Tags: log.Tags,
	}
}
