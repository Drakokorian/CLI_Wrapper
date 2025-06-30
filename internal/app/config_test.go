package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigLoadSave(t *testing.T) {
	dir := t.TempDir()
	wd := t.TempDir()
	cfg := Config{Concurrency: 3, Theme: "dark", CPUThreshold: 50, MemoryThreshold: 50, PollInterval: 2, WorkingDir: wd, LogLevel: "debug", LogPath: filepath.Join(dir, "log.txt")}
	if err := SaveConfig(dir, cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}
	loaded, err := LoadConfig(dir)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if loaded.Concurrency != 3 {
		t.Fatalf("got %d want 3", loaded.Concurrency)
	}
	if loaded.Theme != "dark" {
		t.Fatalf("got %s want dark", loaded.Theme)
	}
	if loaded.CPUThreshold != 50 || loaded.MemoryThreshold != 50 || loaded.WorkingDir != wd || loaded.LogLevel != "debug" || loaded.LogPath == "" {
		t.Fatalf("config not saved")
	}
	// corrupt file should default to 1
	path := filepath.Join(dir, "config", "config.json")
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatalf("write corrupt: %v", err)
	}
	loaded, err = LoadConfig(dir)
	if err == nil {
		t.Fatalf("expected error")
	}

	cfg.Theme = "invalid"
	if err := SaveConfig(dir, cfg); err == nil {
		t.Fatal("expected error for invalid theme")
	}

	cfg.Theme = "dark"
	cfg.CPUThreshold = -1
	if err := SaveConfig(dir, cfg); err == nil {
		t.Fatal("expected error for invalid CPU threshold")
	}
	cfg.CPUThreshold = 50
	cfg.MemoryThreshold = 0
	if err := SaveConfig(dir, cfg); err == nil {
		t.Fatal("expected error for invalid memory threshold")
	}
	cfg.MemoryThreshold = 50
	cfg.PollInterval = 0
	if err := SaveConfig(dir, cfg); err == nil {
		t.Fatal("expected error for invalid poll interval")
	}
	cfg.PollInterval = 2
	cfg.WorkingDir = filepath.Join(dir, "missing")
	if err := SaveConfig(dir, cfg); err == nil {
		t.Fatal("expected error for invalid working dir")
	}

	cfg.WorkingDir = wd
	cfg.LogLevel = "unknown"
	if err := SaveConfig(dir, cfg); err == nil {
		t.Fatal("expected error for invalid log level")
	}
	cfg.LogLevel = "info"
	cfg.LogPath = ""
	if err := SaveConfig(dir, cfg); err == nil {
		t.Fatal("expected error for empty log path")
	}
}
