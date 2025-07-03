package scripts

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestLaunchScriptDryRun(t *testing.T) {
	repoRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	tmp := t.TempDir()
	logPath := filepath.Join(tmp, "sentinel.log")

	// stub wails executable
	wailsPath := filepath.Join(tmp, "wails")
	content := "#!/bin/sh\n\necho wails \"$@\"\n"
	if err := os.WriteFile(wailsPath, []byte(content), 0755); err != nil {
		t.Fatalf("write stub: %v", err)
	}

	cmd := exec.Command("bash", filepath.Join(repoRoot, "scripts", "launch.sh"))
	cmd.Env = append(os.Environ(),
		"PATH="+tmp+":"+os.Getenv("PATH"),
		"SENTINEL_LOG_PATH="+logPath,
		"LAUNCH_SKIP_INSTALL=1",
		"LAUNCH_NOP=1",
	)
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run script: %v output %s", err, out)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	if !bytes.Contains(data, []byte("building")) {
		t.Fatalf("expected build log, got %s", data)
	}
}
