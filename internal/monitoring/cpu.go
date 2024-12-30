package monitoring

import (
	"encoding/json"
	"fmt"

	"github.com/shirou/gopsutil/cpu"
)

type CPUMonitor struct{}

type CPUStats struct {
	Percentage []float64 `json:"percentage"`
}

// GetCpuStats retrieves the current CPU usage percentage since the last call.
func (m *CPUMonitor) GetStats() (interface{}, error) {
	// Retrieve per-CPU times
	percent, err := cpu.Percent(0, true)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CPU times: %v", err)
	}

	// Create a new CPUStats instance
	stats := &CPUStats{
		Percentage: percent,
	}

	return stats, nil
}

func (m *CPUStats) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}
