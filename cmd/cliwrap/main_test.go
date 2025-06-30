package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"cli-wrapper/internal/app"
)

func TestRunCommandUpdatesConfig(t *testing.T) {
	home := t.TempDir()
	os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	work := t.TempDir()
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	cmd := exec.Command("go", "run", "./cmd/cliwrap", "-concurrency", "2", "-workdir", work, "run", "-model", "echo", "-prompt", "hello")
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v output %s", err, out)
	}
	if !strings.Contains(string(out), "hello") {
		t.Fatalf("unexpected output: %s", out)
	}
	base, err := app.AppDir()
	if err != nil {
		t.Fatalf("app dir: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(base, "config", "config.json"))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	var cfg app.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if cfg.Concurrency != 2 || cfg.WorkingDir != work {
		t.Fatalf("config not updated: %+v", cfg)
	}
}

func TestHistoryExportImport(t *testing.T) {
	home1 := t.TempDir()
	os.Setenv("HOME", home1)
	defer os.Unsetenv("HOME")

	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	cmdRun := exec.Command("go", "run", "./cmd/cliwrap", "run", "-model", "echo", "-prompt", "one")
	cmdRun.Dir = repoRoot
	cmdRun.Run()
	exportCmd := exec.Command("go", "run", "./cmd/cliwrap", "history", "export")
	exportCmd.Dir = repoRoot
	exp, err := exportCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("export: %v", err)
	}
	if !bytes.Contains(exp, []byte("one")) {
		t.Fatalf("missing export data: %s", exp)
	}

	home2 := t.TempDir()
	os.Setenv("HOME", home2)
	cmd := exec.Command("go", "run", "./cmd/cliwrap", "history", "import")
	cmd.Dir = repoRoot
	cmd.Stdin = bytes.NewReader(exp)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("import: %v output %s", err, out)
	}
	verifyCmd := exec.Command("go", "run", "./cmd/cliwrap", "history", "export")
	verifyCmd.Dir = repoRoot
	out, err := verifyCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("export2: %v", err)
	}
	if !bytes.Contains(out, []byte("one")) {
		t.Fatalf("imported record missing: %s", out)
	}
}
