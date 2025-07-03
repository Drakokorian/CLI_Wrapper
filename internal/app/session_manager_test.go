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
	m := NewSessionManager(logger, 1, cfg, hist)
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
	m := NewSessionManager(logger, 1, cfg, hist)
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
	m := NewSessionManager(logger, 1, cfg, hist)
	defer m.Close()

	if _, err := m.AddSession("sh", "", []string{"-c", "echo test"}); err != nil {
		t.Fatalf("add: %v", err)
	}
	for len(m.Sessions()) > 0 {
		time.Sleep(10 * time.Millisecond)
	}
}

func TestSessionStreaming(t *testing.T) {
	logger, _ := logging.NewWithPath("info", filepath.Join(t.TempDir(), "log.txt"))
	defer logger.Close()
	cfg := &Config{Concurrency: 1, CPUThreshold: 50, MemoryThreshold: 50, PollInterval: 1, WorkingDir: t.TempDir(), LogLevel: "info", LogPath: filepath.Join(t.TempDir(), "log.txt")}
	base := t.TempDir()
	os.MkdirAll(filepath.Join(base, "state"), 0o755)
	hist, _ := history.New(base)
	defer hist.Close()
	m := NewSessionManager(logger, 1, cfg, hist)
	defer m.Close()

	script := filepath.Join(t.TempDir(), "script.sh")
	os.WriteFile(script, []byte("echo foo\necho bar\n"), 0o755)
	id, err := m.AddSession("sh", "", []string{script})
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	ch, err := m.OutputChannel(id)
	if err != nil {
		t.Fatalf("channel: %v", err)
	}
	var lines []string
	for line := range ch {
		lines = append(lines, line)
	}
	if len(lines) != 2 || lines[0] != "foo" || lines[1] != "bar" {
		t.Fatalf("unexpected lines %v", lines)
	}
}
