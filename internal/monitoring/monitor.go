package monitoring

import (
	"encoding/json"
	"time"
)

// Monitor defines a struct that holds a slice of Monitorable interfaces.
type Monitor struct {
	monitors []Monitorable
}

// Monitorable is an interface that requires a GetStats method.
type Monitorable interface {
	GetStats() (interface{}, error)
}

// Stats aggregates the statistics from all monitors.
type Stats struct {
	Time    time.Time     `json:"time"`
	CPU     *CPUStats     `json:"cpu"`
	Mem     *MemStats     `json:"mem"`
	Network *NetworkStats `json:"network"`
}

// NewMonitor initializes a new Monitor with the provided Monitorable instances.
func NewMonitor(monitors ...Monitorable) *Monitor {
	return &Monitor{
		monitors: monitors,
	}
}

// GetStats retrieves the current statistics from all monitors.
func (m *Monitor) GetStats() (Stats, error) {
	stats := Stats{
		Time: time.Now(),
	}

	for _, monitor := range m.monitors {
		stat, err := monitor.GetStats()
		if err != nil {
			return Stats{}, err
		}

		switch monitor.(type) {
		case *CPUMonitor:
			stats.CPU = stat.(*CPUStats)
		case *MemMonitor:
			stats.Mem = stat.(*MemStats)
		case *NetworkMonitor:
			stats.Network = stat.(*NetworkStats)
		}
	}

	return stats, nil
}

func (st *Stats) String() string {
	s, _ := json.Marshal(st)
	return string(s)
}
