package logging

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"
)

type rotateWriter struct {
	path       string
	file       *os.File
	maxSize    int64
	maxBackups int
	maxAge     time.Duration
}

func newRotateWriter(path string, maxSize int64, maxBackups int, maxAge time.Duration) (*rotateWriter, error) {
	if maxSize <= 0 || maxBackups < 0 {
		return nil, errors.New("invalid rotation parameters")
	}
	rw := &rotateWriter{
		path:       path,
		maxSize:    maxSize,
		maxBackups: maxBackups,
		maxAge:     maxAge,
	}
	if err := rw.open(); err != nil {
		return nil, err
	}
	return rw, nil
}

func (w *rotateWriter) open() error {
	dir := filepath.Dir(w.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create log dir: %w", err)
	}
	f, err := os.OpenFile(w.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}
	w.file = f
	return nil
}

func (w *rotateWriter) Write(p []byte) (int, error) {
	info, err := w.file.Stat()
	if err != nil {
		return 0, fmt.Errorf("stat log file: %w", err)
	}
	if info.Size()+int64(len(p)) > w.maxSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}
	return w.file.Write(p)
}

func (w *rotateWriter) rotate() error {
	if err := w.file.Close(); err != nil {
		return fmt.Errorf("close log file: %w", err)
	}
	ts := time.Now().UTC().Format("20060102T150405Z")
	rotated := fmt.Sprintf("%s-%s", w.path, ts)
	if err := os.Rename(w.path, rotated); err != nil {
		return fmt.Errorf("rotate log file: %w", err)
	}
	if err := w.cleanup(); err != nil {
		return err
	}
	return w.open()
}

func (w *rotateWriter) cleanup() error {
	files, err := filepath.Glob(w.path + "-*")
	if err != nil {
		return fmt.Errorf("glob backups: %w", err)
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i] > files[j]
	})
	now := time.Now().UTC()
	valid := files[:0]
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			continue
		}
		if now.Sub(info.ModTime()) > w.maxAge {
			os.Remove(f)
			continue
		}
		valid = append(valid, f)
	}
	if len(valid) > w.maxBackups {
		remove := valid[w.maxBackups:]
		for _, f := range remove {
			os.Remove(f)
		}
	}
	return nil
}

func (w *rotateWriter) Close() error {
	return w.file.Close()
}

type Logger struct {
	writer io.WriteCloser
}

type logEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

func New() (*Logger, error) {
	path := "/var/log/sentinel.log"
	if runtime.GOOS == "windows" {
		path = filepath.Join(os.Getenv("SystemDrive"), "Temp", "sentinel.log")
	}
	rw, err := newRotateWriter(path, 20*1024*1024, 5, 30*24*time.Hour)
	if err != nil {
		return nil, err
	}
	return &Logger{writer: rw}, nil
}

func (l *Logger) log(level, msg string) error {
	entry := logEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   msg,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal log: %w", err)
	}
	if _, err := l.writer.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("write log: %w", err)
	}
	return nil
}

func (l *Logger) Info(msg string) error {
	return l.log("INFO", msg)
}

func (l *Logger) Error(msg string) error {
	return l.log("ERROR", msg)
}

func (l *Logger) Close() error {
	return l.writer.Close()
}
