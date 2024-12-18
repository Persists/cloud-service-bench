package generator

import (
	"cloud-service-bench/internal/connection"
	"cloud-service-bench/internal/log"
	"sync"
	"time"
)

// routine is a function that sends logs to a connection client
// therefore it prepares the logs by adding a timestamp and sends them
func routine(
	fluentdConfig *FluentdConfig,
	logChan chan *log.LogMessage,
	ready *sync.WaitGroup,
) {
	connectionClient := connection.NewConnectionClient(fluentdConfig.Host, fluentdConfig.Port)
	connectionClient.Connect()
	defer connectionClient.Disconnect()

	ready.Done()
	ready.Wait()

	for logMessage := range logChan {
		logMessage.Time = time.Now()
		err := connectionClient.SendMessage(logMessage.ToFluentdMessage())
		if err != nil {
			panic(err)
		}
	}
}
