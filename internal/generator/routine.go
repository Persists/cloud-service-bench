package generator

import (
	"cloud-service-bench/internal/archive"
	"cloud-service-bench/internal/connection"
	"cloud-service-bench/internal/log"
	"fmt"
	"sync"
	"time"

	"cloud-service-bench/internal/config"
)

func routine(
	TcpConfig *config.TcpConfig,
	samples []*log.LogMessage,
	ticksPerSecond float64,
	messagesPerTick int,
	ready *sync.WaitGroup,
	stop chan struct{},
	archiver archive.Archiver,
) error {
	connectionClient := connection.NewConnectionClient(TcpConfig.Host, TcpConfig.Port)
	err := connectionClient.Connect()
	if err != nil {
		return err
	}

	defer connectionClient.Disconnect()

	ready.Done()
	ready.Wait()

	ticker := time.NewTicker(time.Second / time.Duration(ticksPerSecond))
	index := 0

	for {
		select {
		case <-ticker.C:
			for i := 0; i < messagesPerTick; i++ {
				log := *samples[index]
				log.Time = time.Now()
				err = connectionClient.SendMessage(log.ToFluentdMessage())
				if err != nil {
					fmt.Errorf("failed to send message: %w", err)
				}
				index = (index + 1) % len(samples)
				archiver.Write(log.ToFluentdMessage())
			}
		case <-stop:
			return nil
		}
	}
}
