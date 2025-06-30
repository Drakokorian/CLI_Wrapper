package app

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"cli-wrapper/internal/history"
	"cli-wrapper/internal/logging"
)

func TestSessionManagerQueue(t *testing.T) {
	logger, _ := logging.NewWithPath(filepath.Join(t.TempDir(), "log.txt"))
	defer logger.Close()
	cfg := &Config{Concurrency: 1, CPUThreshold: 50, MemoryThreshold: 50, PollInterval: 1}
	hist, _ := history.New(t.TempDir())
	defer hist.Close()
	m := NewSessionManager(t.TempDir(), logger, 1, cfg, hist)
	defer m.Close()

	start := time.Now()
	id1, err := m.AddSession("sh", "", []string{"-c", "sleep 0.1"})
	if err != nil {
		t.Fatalf("add1: %v", err)
	}
	id2, err := m.AddSession("sh", "", []string{"-c", "sleep 0.1"})
	if err != nil {
		t.Fatalf("add2: %v", err)
	}
	if id1 == id2 {
		t.Fatalf("duplicate ids")
	}
	for len(m.Sessions()) > 0 {
		time.Sleep(10 * time.Millisecond)
	}
	if time.Since(start) < 190*time.Millisecond {
		t.Fatalf("sessions did not run sequentially")
	}
}

func TestSessionTerminate(t *testing.T) {
	logger, _ := logging.NewWithPath(filepath.Join(t.TempDir(), "log.txt"))
	defer logger.Close()
	cfg := &Config{Concurrency: 1, CPUThreshold: 50, MemoryThreshold: 50, PollInterval: 1}
	hist, _ := history.New(t.TempDir())
	defer hist.Close()
	m := NewSessionManager(t.TempDir(), logger, 1, cfg, hist)
	defer m.Close()

	id, err := m.AddSession("sh", "", []string{"-c", "sleep 2"})
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	time.Sleep(50 * time.Millisecond)
	if err := m.Terminate(id); err != nil {
		t.Fatalf("terminate: %v", err)
	}
}

func TestSessionWorkingDir(t *testing.T) {
	logger, _ := logging.NewWithPath(filepath.Join(t.TempDir(), "log.txt"))
	defer logger.Close()
	work := t.TempDir()
	cfg := &Config{Concurrency: 1, CPUThreshold: 50, MemoryThreshold: 50, PollInterval: 1, WorkingDir: work}
	hist, _ := history.New(t.TempDir())
	defer hist.Close()
	m := NewSessionManager(t.TempDir(), logger, 1, cfg, hist)
	defer m.Close()

	outFile := filepath.Join(t.TempDir(), "out.txt")
	cmd := "pwd > " + outFile
	if _, err := m.AddSession("sh", "", []string{"-c", cmd}); err != nil {
		t.Fatalf("add: %v", err)
	}
	for len(m.Sessions()) > 0 {
		time.Sleep(10 * time.Millisecond)
	}
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("read out: %v", err)
	}
	got := string(bytes.TrimSpace(data))
	if got != work {
		t.Fatalf("got %s want %s", got, work)
	}
}
