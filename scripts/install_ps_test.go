package scripts

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestInstallPSLog(t *testing.T) {
	pwsh, err := exec.LookPath("pwsh")
	if err != nil {
		t.Skip("pwsh not installed")
	}
	repoRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	tmp := t.TempDir()
	logPath := filepath.Join(tmp, "sentinel.log")
	cmd := exec.Command(pwsh, filepath.Join(repoRoot, "scripts", "install.ps1"))
	cmd.Env = append(os.Environ(),
		"SENTINEL_LOG_PATH="+logPath,
		"INSTALL_NOP=1",
	)
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run script: %v output %s", err, out)
	}
	if _, err := os.Stat(logPath); err != nil {
		t.Fatalf("log not created: %v", err)
	}
}
