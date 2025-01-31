package monitoring

import (
	"encoding/json"

	"github.com/shirou/gopsutil/mem"
)

type MemMonitor struct{}

type MemStats struct {
	Total      uint64 `json:"total"`
	Used       uint64 `json:"used"`
	Free       uint64 `json:"free"`
	SwapFree   uint64 `json:"swap_free"`
	SwapTotal  uint64 `json:"swap_total"`
	SwapCached uint64 `json:"swap_cached"`
	Inactive   uint64 `json:"inactive"`
	Cached     uint64 `json:"cached"`
}

// GetMemStats retrieves the current memory usage statistics.
func (m *MemMonitor) GetStats() (interface{}, error) {
	mem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	stats := &MemStats{
		Total:      mem.Total,
		Used:       mem.Used,
		Free:       mem.Free,
		Cached:     mem.Cached,
		SwapFree:   mem.SwapFree,
		SwapTotal:  mem.SwapTotal,
		SwapCached: mem.SwapCached,
		Inactive:   mem.Inactive,
	}

	return stats, nil
}

func (m *MemStats) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}
