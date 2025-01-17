package generator

import (
	"cloud-service-bench/internal/config"
	"cloud-service-bench/internal/connection"
	"cloud-service-bench/internal/log"
	"fmt"
	"time"
)

type Worker struct {
	ID               int
	connectionClient *connection.ConnectionClient
	stop             chan struct{}
	start            chan struct{}
	TotalMessages    int
}

func newWorker(tcpConfig *config.TcpConfig, sampleLogs []*log.LogMessage, startChan chan struct{}, stopChan chan struct{}, id int) (*Worker, error) {
	cc := connection.NewConnectionClient(
		tcpConfig.Host,
		tcpConfig.Port,
	)

	err := cc.Connect()
	if err != nil {
		return nil, err
	}

	worker := &Worker{
		ID:               id,
		connectionClient: cc,
		stop:             stopChan,
		start:            startChan,
		TotalMessages:    0,
	}

	go worker.routine(sampleLogs)

	return worker, nil
}

func (w *Worker) routine(
	samples []*log.LogMessage,
) {
	defer w.connectionClient.Disconnect()
	<-w.start
	for {
		select {
		case <-w.stop:
			fmt.Printf("worker %d stopped\n", w.ID)
			fmt.Printf("total messages sent by worker %d: %d\n", w.ID, w.TotalMessages)
			return
		default:
			log := *samples[w.TotalMessages%len(samples)]
			log.Tags = []string{fmt.Sprintf("worker-%d", w.ID), fmt.Sprintf("log-%d", w.TotalMessages)}
			log.Time = log.Time.Add(time.Duration(w.TotalMessages) * time.Second)

			err := w.connectionClient.SendMessage(log.ToFluentdMessage())
			if err != nil {
				fmt.Printf("failed to send message: %v\n", err)
				break
			}
			w.TotalMessages++
		}
	}
}
