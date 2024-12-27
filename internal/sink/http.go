package sink

import (
	"cloud-service-bench/internal/archive"
	"cloud-service-bench/internal/log"
	"fmt"
	"io"
	"net/http"
	"time"
)

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

	logEntryCount := len(logMessages)

	for _, logMessage := range logMessages {
		archivable := fmt.Sprintf("%s: [%s]", requestTime.Format("2006-01-02T15:04:05.000Z"), logMessage.ToArchivable())
		hs.archiver.Write(archivable)
	}

	fmt.Println("Request processed in", time.Since(requestTime))
	fmt.Println("Archived", logEntryCount, "log entries")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Request received successfully")
}
