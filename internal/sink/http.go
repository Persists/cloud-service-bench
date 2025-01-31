package sink

import (
	"cloud-service-bench/internal/archive"
	"cloud-service-bench/internal/log"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LogBatch struct {
	ArrivalTime time.Time            `json:"arrivalTime"`
	LogMessages []log.SlimLogMessage `json:"logMessages"`
}

type HttpSink struct {
	archiver archive.Archiver
}

func NewHttpSink(archiver archive.Archiver) *HttpSink {
	return &HttpSink{
		archiver: archiver,
	}
}

// Handler is the HTTP handler for the sink
// It reads the body of the request, decodes the JSON log messages, and writes them to the archiver
func (hs *HttpSink) Handler(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now()
	fmt.Println("Received request at", requestTime)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}

	logMessages, err := log.DecodeJson(body)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// to reduce the size of the log messages before archiving
	var slimLogMessages []log.SlimLogMessage
	for _, logMessage := range logMessages {
		slimLogMessages = append(slimLogMessages, logMessage.ToSlimLogMessage())
	}
	fmt.Println("Received", len(logMessages), "log messages")

	logBatch := LogBatch{
		ArrivalTime: requestTime,
		LogMessages: slimLogMessages,
	}

	// writing the log batch to the archiver
	// if there is an error, write an error message to the archiver
	jsonLogBatch, err := json.Marshal(logBatch)
	if err != nil {
		hs.archiver.Write(fmt.Sprintf("%s | Error marshalling | Batch length: %d", requestTime.Format("2006-01-02T15:04:05.000Z"), len(logMessages)))
		return
	}
	hs.archiver.Write(string(jsonLogBatch))

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Request received successfully")
}
