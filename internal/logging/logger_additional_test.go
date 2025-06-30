package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewRotateWriterInvalid(t *testing.T) {
	if _, err := newRotateWriter("", 0, -1, time.Hour); err == nil {
		t.Fatal("expected error")
	}
}

func TestNewRotateWriterOpenError(t *testing.T) {
	if _, err := newRotateWriter(filepath.Join("/proc", "noexist", "x.log"), 10, 1, time.Hour); err == nil {
		t.Fatal("expected error")
	}
}

func TestRotateAndCleanup(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sentinel.log")
	rw, err := newRotateWriter(path, 20, 2, time.Hour)
	if err != nil {
		t.Fatalf("newRotateWriter: %v", err)
	}
	defer rw.Close()

	old := path + "-old"
	if err := os.WriteFile(old, []byte("old"), 0o644); err != nil {
		t.Fatalf("write old: %v", err)
	}
	oldTime := time.Now().Add(-2 * time.Hour)
	if err := os.Chtimes(old, oldTime, oldTime); err != nil {
		t.Fatalf("chtimes: %v", err)
	}

	for i := 0; i < 4; i++ {
		if _, err := rw.Write([]byte(strings.Repeat("x", 15))); err != nil {
			t.Fatalf("write: %v", err)
		}
		time.Sleep(time.Millisecond * 1100)
	}

	backups, err := filepath.Glob(path + "-*")
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(backups) != 2 {
		t.Fatalf("expected 2 backups got %d", len(backups))
	}
	for _, f := range backups {
		info, _ := os.Stat(f)
		if info.ModTime().Before(oldTime) {
			t.Fatal("old backup not cleaned up")
		}
	}
}

func TestLoggerLevels(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sentinel.log")
	logger, err := NewWithPath("debug", path)
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}
	defer logger.Close()
	if err := logger.Debug("dbg"); err != nil {
		t.Fatalf("debug: %v", err)
	}
	if err := logger.Info("hello"); err != nil {
		t.Fatalf("info: %v", err)
	}
	if err := logger.Error("bad"); err != nil {
		t.Fatalf("error: %v", err)
	}
	if err := logger.ModelSwitch("a", "b", "prompt"); err != nil {
		t.Fatalf("switch: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines got %d", len(lines))
	}
	for _, l := range lines {
		var obj map[string]any
		if err := json.Unmarshal([]byte(l), &obj); err != nil {
			t.Fatalf("json: %v", err)
		}
	}
}

type errWriter struct{}

func (errWriter) Write([]byte) (int, error) { return 0, fmt.Errorf("write fail") }
func (errWriter) Close() error              { return nil }

func TestLogWriteError(t *testing.T) {
	l := &Logger{writer: errWriter{}, level: LevelInfo}
	if err := l.Info("bad"); err == nil {
		t.Fatal("expected error")
	}
}

func TestModelSwitchWriteError(t *testing.T) {
	l := &Logger{writer: errWriter{}, level: LevelInfo}
	if err := l.ModelSwitch("a", "b", "p"); err == nil {
		t.Fatal("expected error")
	}
}

func TestWriteStatError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sentinel.log")
	rw, err := newRotateWriter(path, 20, 1, time.Hour)
	if err != nil {
		t.Fatalf("newRotateWriter: %v", err)
	}
	rw.file.Close()
	if _, err := rw.Write([]byte("x")); err == nil {
		t.Fatal("expected error")
	}
}

func TestRotateOpenError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sentinel.log")
	rw, err := newRotateWriter(path, 5, 1, time.Hour)
	if err != nil {
		t.Fatalf("newRotateWriter: %v", err)
	}
	if _, err := rw.Write([]byte("abcdef")); err != nil {
		t.Fatalf("write: %v", err)
	}
	rw.path = filepath.Join("/proc", "noexist", "log")
	if _, err := rw.Write([]byte("more")); err == nil {
		t.Fatal("expected error")
	}
}

func TestNewDefault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sentinel.log")
	l, err := NewWithPath("info", path)
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	if err := l.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}

func TestRotateWriterMultiple(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sentinel.log")
	rw, err := newRotateWriter(path, 10, 1, time.Hour)
	if err != nil {
		t.Fatalf("newRotateWriter: %v", err)
	}
	defer rw.Close()
	for i := 0; i < 3; i++ {
		if _, err := rw.Write([]byte("0123456789")); err != nil {
			t.Fatalf("write: %v", err)
		}
		time.Sleep(time.Millisecond * 1100)
	}
	backups, err := filepath.Glob(path + "-*")
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(backups) != 1 {
		t.Fatalf("expected 1 backup got %d", len(backups))
	}
}

func TestRotateManual(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sentinel.log")
	rw, err := newRotateWriter(path, 20, 1, time.Hour)
	if err != nil {
		t.Fatalf("newRotateWriter: %v", err)
	}
	if _, err := rw.file.Write([]byte("hello")); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := rw.rotate(); err != nil {
		t.Fatalf("rotate: %v", err)
	}
	backups, err := filepath.Glob(path + "-*")
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(backups) != 1 {
		t.Fatalf("expected 1 backup, got %d", len(backups))
	}
	rw.Close()
}
