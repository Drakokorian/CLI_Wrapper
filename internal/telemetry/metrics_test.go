package telemetry

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

func TestRead(t *testing.T) {
	u, err := Read()
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if u.CPU < 0 || u.CPU > 100 {
		t.Fatalf("cpu out of range: %v", u.CPU)
	}
	if u.Memory < 0 || u.Memory > 100 {
		t.Fatalf("mem out of range: %v", u.Memory)
	}
}

func TestReadProcess(t *testing.T) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		t.Fatalf("new process: %v", err)
	}
	u, err := ReadProcess(p)
	if err != nil {
		t.Fatalf("read process: %v", err)
	}
	if u.CPU < 0 || u.CPU > 100 {
		t.Fatalf("cpu out of range: %v", u.CPU)
	}
	if u.Memory < 0 || u.Memory > 100 {
		t.Fatalf("mem out of range: %v", u.Memory)
	}
}

// fdCount returns the number of open file descriptors for the current process.
func fdCount() (int, error) {
	entries, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		return 0, err
	}
	return len(entries), nil
}

func TestReadProcess_NoLeak(t *testing.T) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		t.Fatalf("new process: %v", err)
	}
	start, err := fdCount()
	if err != nil {
		t.Skipf("fd count unsupported: %v", err)
	}
	for i := 0; i < 5; i++ {
		if _, err := ReadProcess(p); err != nil {
			t.Fatalf("read process: %v", err)
		}
	}
	end, err := fdCount()
	if err != nil {
		t.Fatalf("fd count: %v", err)
	}
	if end > start {
		t.Fatalf("fd leak: start %d end %d", start, end)
	}
}

func TestPollProcess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch, err := PollProcess(ctx, int32(os.Getpid()), 10*time.Millisecond)
	if err != nil {
		t.Fatalf("poll: %v", err)
	}
	select {
	case u := <-ch:
		if u.CPU < 0 || u.CPU > 100 || u.Memory < 0 || u.Memory > 100 {
			t.Fatalf("invalid usage: %#v", u)
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatal("timeout")
	}
	cancel()
	for range ch {
	}
}
