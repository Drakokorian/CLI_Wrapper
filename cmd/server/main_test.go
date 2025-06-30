package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cli-wrapper/internal/app"
	"cli-wrapper/internal/logging"
)

func TestServerStartupLogFailure(t *testing.T) {
	home := t.TempDir()
	os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	base, err := app.PrepareDirectories()
	if err != nil {
		t.Fatalf("prep dirs: %v", err)
	}
	cfgDir := filepath.Join(base, "config")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("mk config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, "config.json"), []byte("bad"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	logPath := filepath.Join(t.TempDir(), "log.txt")
	logger, err := logging.NewWithPath("info", logPath)
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}
	defer logger.Close()

	_, err = app.LoadConfig(base)
	if err != nil {
		logger.Error("load config: " + err.Error())
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	if !strings.Contains(string(data), "load config") {
		t.Fatalf("log missing: %s", string(data))
	}
}
