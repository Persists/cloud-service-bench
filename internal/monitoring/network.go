package monitoring

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
)

type NetworkMonitor struct{}

type NetworkStats struct {
	net.IOCountersStat
}

// GetNetworkStats retrieves the current network I/O statistics.
func (m *NetworkMonitor) GetStats() (interface{}, error) {
	ioCounters, err := net.IOCounters(false)
	if err != nil {
		return nil, fmt.Errorf("error retrieving network I/O statistics: %v", err)
	}

	stats := &NetworkStats{
		ioCounters[0],
	}

	return stats, nil
}
