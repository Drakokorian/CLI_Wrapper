package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectCLIToolNone(t *testing.T) {
	t.Setenv("CLIWRAP_TOOL", "")
	if _, err := DetectCLITool(); err == nil {
		t.Fatal("expected error")
	}
}

func TestDetectCLIToolsNone(t *testing.T) {
	t.Setenv("CLIWRAP_TOOL", "")
	if _, err := DetectCLITools(); err == nil {
		t.Fatal("expected error")
	}
}

func TestDetectCLIToolEnv(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "openai")
	if err := os.WriteFile(script, []byte("#!/bin/sh\necho ok"), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}
	oldPath := os.Getenv("PATH")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+oldPath)
	t.Setenv("CLIWRAP_TOOL", "openai")

	tool, err := DetectCLITool()
	if err != nil {
		t.Fatalf("detect: %v", err)
	}
	if tool != "openai" {
		t.Fatalf("got %s want openai", tool)
	}
}
