package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
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

// ReadProcess returns CPU and memory usage for the provided process.
// It expects a valid gopsutil process handle.
func ReadProcess(p *process.Process) (Usage, error) {
	if p == nil {
		return Usage{}, fmt.Errorf("nil process")
	}
	cpuPercent, err := p.CPUPercent()
	if err != nil {
		return Usage{}, fmt.Errorf("cpu percent: %w", err)
	}
	memPercent, err := p.MemoryPercent()
	if err != nil {
		return Usage{}, fmt.Errorf("memory percent: %w", err)
	}
	return Usage{CPU: cpuPercent, Memory: float64(memPercent)}, nil
}

// PollProcess periodically emits Usage metrics for the specified process ID.
// The returned channel is closed when ctx is done.
func PollProcess(ctx context.Context, pid int32, interval time.Duration) (<-chan Usage, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("invalid interval")
	}
	p, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("new process: %w", err)
	}
	ch := make(chan Usage)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			case <-ticker.C:
				cpuPercent, err := p.CPUPercent()
				if err != nil {
					continue
				}
				memPercent, err := p.MemoryPercent()
				if err != nil {
					continue
				}
				ch <- Usage{CPU: cpuPercent, Memory: float64(memPercent)}
			}
		}
	}()
	return ch, nil
}
