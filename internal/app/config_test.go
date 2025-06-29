package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigLoadSave(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{Concurrency: 3, Theme: "dark"}
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
}
