//go:build windows

package app

import (
	"os/exec"
	"testing"
)

// Test that throttleProcess returns an error when the process has not been started.
func TestThrottleProcessNotStarted(t *testing.T) {
	cmd := exec.Command("cmd.exe")
	if err := throttleProcess(cmd); err == nil {
		t.Fatal("expected error")
	}
}

// Test that throttleProcess propagates errors from OpenProcess.
func TestThrottleProcessOpenError(t *testing.T) {
	cmd := exec.Command("cmd.exe", "/C", "exit 0")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	// Wait for the process to exit so OpenProcess fails.
	if err := cmd.Wait(); err != nil {
		t.Fatalf("wait: %v", err)
	}
	if err := throttleProcess(cmd); err == nil {
		t.Fatal("expected error")
	}
}
