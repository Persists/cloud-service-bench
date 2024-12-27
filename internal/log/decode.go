package log

import "encoding/json"

func DecodeJson(body []byte) ([]LogMessage, error) {
	var logMessages []LogMessage
	if err := json.Unmarshal(body, &logMessages); err != nil {
		return nil, err
	}
	return logMessages, nil
}
