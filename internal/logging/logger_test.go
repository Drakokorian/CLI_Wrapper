package logging

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoggerWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sentinel.log")
	rw, err := newRotateWriter(path, 1024, 1, time.Hour)
	if err != nil {
		t.Fatalf("newRotateWriter: %v", err)
	}
	logger := &Logger{writer: rw}
	defer logger.Close()
	if err := logger.Info("test"); err != nil {
		t.Fatalf("info: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected data written")
	}
}
