package app

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"cli-wrapper/internal/history"
	"cli-wrapper/internal/logging"
)

func TestSessionManagerQueue(t *testing.T) {
	logger, _ := logging.NewWithPath("info", filepath.Join(t.TempDir(), "log.txt"))
	defer logger.Close()
	cfg := &Config{Concurrency: 1, CPUThreshold: 50, MemoryThreshold: 50, PollInterval: 1, LogLevel: "info", LogPath: filepath.Join(t.TempDir(), "log.txt")}
	base := t.TempDir()
	os.MkdirAll(filepath.Join(base, "state"), 0o755)
	hist, _ := history.New(base)
	defer hist.Close()
	m := NewSessionManager(t.TempDir(), logger, 1, cfg, hist)
	defer m.Close()

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
}

func TestSessionTerminate(t *testing.T) {
	logger, _ := logging.NewWithPath("info", filepath.Join(t.TempDir(), "log.txt"))
	defer logger.Close()
	cfg := &Config{Concurrency: 1, CPUThreshold: 50, MemoryThreshold: 50, PollInterval: 1, LogLevel: "info", LogPath: filepath.Join(t.TempDir(), "log.txt")}
	base := t.TempDir()
	os.MkdirAll(filepath.Join(base, "state"), 0o755)
	hist, _ := history.New(base)
	defer hist.Close()
	m := NewSessionManager(t.TempDir(), logger, 1, cfg, hist)
	defer m.Close()

	id, err := m.AddSession("sh", "", []string{"-c", "sleep 2"})
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	time.Sleep(50 * time.Millisecond)
	_ = m.Terminate(id)
}

func TestSessionWorkingDir(t *testing.T) {
	logger, _ := logging.NewWithPath("info", filepath.Join(t.TempDir(), "log.txt"))
	defer logger.Close()
	work := t.TempDir()
	cfg := &Config{Concurrency: 1, CPUThreshold: 50, MemoryThreshold: 50, PollInterval: 1, WorkingDir: work, LogLevel: "info", LogPath: filepath.Join(t.TempDir(), "log.txt")}
	base := t.TempDir()
	os.MkdirAll(filepath.Join(base, "state"), 0o755)
	hist, _ := history.New(base)
	defer hist.Close()
	m := NewSessionManager(t.TempDir(), logger, 1, cfg, hist)
	defer m.Close()

	if _, err := m.AddSession("sh", "", []string{"-c", "echo test"}); err != nil {
		t.Fatalf("add: %v", err)
	}
	for len(m.Sessions()) > 0 {
		time.Sleep(10 * time.Millisecond)
	}
}
