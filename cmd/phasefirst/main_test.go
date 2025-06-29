package main

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestPhasefirstMain(t *testing.T) {
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	cases := []struct {
		name    string
		dir     string
		wantErr bool
	}{
		{"success", repoRoot, false},
		{"missingDocs", t.TempDir(), true},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", "./cmd/phasefirst")
			cmd.Dir = tt.dir
			out, err := cmd.CombinedOutput()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil with output %s", out)
				}
				return
			}
			if err != nil {
				t.Fatalf("run err: %v output %s", err, out)
			}
			var m map[string]string
			if derr := json.Unmarshal(out, &m); derr != nil {
				t.Fatalf("unmarshal: %v output: %s", derr, out)
			}
			if _, ok := m["phase01_core_app_foundation.md"]; !ok {
				t.Fatalf("missing phase01 entry in %v", m)
			}
		})
	}
}
