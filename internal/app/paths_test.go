package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAppDirForOS(t *testing.T) {
	oldAppData := os.Getenv("APPDATA")
	os.Setenv("APPDATA", `C:\\Data`)
	defer os.Setenv("APPDATA", oldAppData)
	path, err := appDirForOS("windows")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filepath.Base(path) != "ai-cli-ui" {
		t.Fatalf("unexpected base: %q", filepath.Base(path))
	}

	home := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", home)
	defer os.Setenv("HOME", oldHome)
	path, err = appDirForOS("darwin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := filepath.Join(home, "Library", "Application Support", "ai-cli-ui")
	if path != expected {
		t.Fatalf("got %q want %q", path, expected)
	}
}
