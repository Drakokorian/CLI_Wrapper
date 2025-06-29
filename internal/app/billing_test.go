package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestBillingURL(t *testing.T) {
	url, err := BillingURL("openai")
	if err != nil || url == "" {
		t.Fatalf("openai billing url err=%v url=%s", err, url)
	}
	url, err = BillingURL("gemini")
	if err != nil || url == "" {
		t.Fatalf("gemini billing url err=%v url=%s", err, url)
	}
	if _, err := BillingURL("foo"); err == nil {
		t.Fatal("expected error for unknown tool")
	}
}

func TestFetchUsage(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "openai")
	content := "#!/bin/sh\necho '{\"total_granted\":100,\"total_used\":50,\"total_available\":50}'"
	if err := os.WriteFile(script, []byte(content), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", dir+string(os.PathListSeparator)+oldPath)
	defer os.Setenv("PATH", oldPath)

	u, err := FetchUsage("openai")
	if err != nil {
		t.Fatalf("fetch usage: %v", err)
	}
	if u.TotalGranted != 100 || u.TotalUsed != 50 || u.TotalAvailable != 50 {
		b, _ := json.Marshal(u)
		t.Fatalf("unexpected usage %s", b)
	}
	if _, err := FetchUsage("unknown"); err == nil {
		t.Fatal("expected error for unknown tool")
	}
}
