package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"cli-wrapper/internal/logging"
)

func TestRunSuccess(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "log.txt")
	logger, err := logging.NewWithPath(logPath)
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}
	defer logger.Close()

	if err := run(logger, "go", "version"); err != nil {
		t.Fatalf("run: %v", err)
	}
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	if !bytes.Contains(data, []byte("go:")) {
		t.Fatalf("missing log output: %s", data)
	}
}

func TestRunFailure(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "log.txt")
	logger, err := logging.NewWithPath(logPath)
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}
	defer logger.Close()

	if err := run(logger, "go", "unknowncmd"); err == nil {
		t.Fatalf("expected error")
	}
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	if !bytes.Contains(data, []byte("failed")) {
		t.Fatalf("missing failure log: %s", data)
	}
}
