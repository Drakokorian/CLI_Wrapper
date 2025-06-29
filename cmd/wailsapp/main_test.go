package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cli-wrapper/internal/app"
)

func TestAppStartupLogFailure(t *testing.T) {
	home := t.TempDir()
	os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	base, err := app.AppDir()
	if err != nil {
		t.Fatalf("app dir: %v", err)
	}
	cfgDir := filepath.Join(base, "config")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("mk config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, "config.json"), []byte("bad"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	main()

	logPath := filepath.Join(base, "logs", "logs.txt")
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	if !strings.Contains(string(data), "load config") {
		t.Fatalf("log missing: %s", string(data))
	}
}
