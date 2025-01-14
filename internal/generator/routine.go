package generator

import (
	"cloud-service-bench/internal/connection"
	"cloud-service-bench/internal/log"
	"fmt"
	"sync"
	"time"
)

// routine is responsible for sending log messages at a certain rate.
// Each routine has a own connection to the server.
// It sends messages at a rate of g.GeneratorConfig.LogsPerSecond.
// It sends messages in batches, the amount of batches per second is defined by g.GeneratorConfig.BatchesPerSec.
// The function waits on a channel to stop sending messages.
func (g *GeneratorClient) routine(
	workerID int,
	samples []*log.LogMessage,
	ready *sync.WaitGroup,
	stop chan struct{},
) error {
	connectionClient := connection.NewConnectionClient(g.TCPConfig.Host, g.TCPConfig.Port)
	err := connectionClient.Connect()
	if err != nil {
		return err
	}

	defer connectionClient.Disconnect()

	ready.Done()
	ready.Wait()

	ticker := time.NewTicker(time.Second / time.Duration(g.GeneratorConfig.BatchesPerSec))
	messagesPerTick := g.GeneratorConfig.LogsPerSecond / g.GeneratorConfig.BatchesPerSec

	globCounter := 0
	for {
		select {
		case <-ticker.C:
			for i := 0; i < messagesPerTick; i++ {
				log := *samples[globCounter%len(samples)]
				log.Time = time.Now()
				// The tags are used to identify the worker and the message number later on.
				log.Tags = []string{fmt.Sprintf("worker-%d", workerID), fmt.Sprintf("%d", globCounter)}
				err = connectionClient.SendMessage(log.ToFluentdMessage())
				if err != nil {
					fmt.Printf("failed to send message: %v\n", err)
					// Break here, to prevent the routine from sending more messages, it will continue with the next tick.
					break
				}
				globCounter++
			}
		case <-stop:
			return nil
		}
	}
}
