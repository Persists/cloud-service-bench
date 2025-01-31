package log

import (
	"encoding/json"
)

// DecodeJson decodes a JSON byte array into a slice of LogMessages.
func DecodeJson(body []byte) ([]LogMessage, error) {
	var logMessages []LogMessage
	if err := json.Unmarshal(body, &logMessages); err != nil {
		return nil, err
	}
	return logMessages, nil
}

// ToSlimLogMessage converts a LogMessage to a SlimLogMessage
// for archiving purposes.
//
// This is used to reduce the size of the log message before.
func (log *LogMessage) ToSlimLogMessage() SlimLogMessage {
	return SlimLogMessage{
		Time: log.Time,
		Name: log.Name,
		Tags: log.Tags,
	}
}
