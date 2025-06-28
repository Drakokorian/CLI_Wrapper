package telemetry

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// Usage represents CPU and memory usage percentages.
type Usage struct {
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
}

// Read returns the current system CPU and memory utilization.
func Read() (Usage, error) {
	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		return Usage{}, fmt.Errorf("cpu percent: %w", err)
	}
	vm, err := mem.VirtualMemory()
	if err != nil {
		return Usage{}, fmt.Errorf("memory percent: %w", err)
	}
	var cpuVal float64
	if len(cpuPercents) > 0 {
		cpuVal = cpuPercents[0]
	}
	return Usage{CPU: cpuVal, Memory: vm.UsedPercent}, nil
}
