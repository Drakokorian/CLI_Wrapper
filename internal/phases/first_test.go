package phases

import (
	"path/filepath"
	"testing"
)

func TestFirstStep(t *testing.T) {
	path := filepath.Join("..", "..", "docs", "phases", "phase01_core_app_foundation.md")
	step, err := FirstStep(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "Scaffold a new Wails v2 project using a Go backend and React frontend."
	if step != expected {
		t.Fatalf("got %q want %q", step, expected)
	}
}

func TestPhaseFirstSteps(t *testing.T) {
	dir := filepath.Join("..", "..", "docs", "phases")
	steps, err := PhaseFirstSteps(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(steps) == 0 {
		t.Fatal("expected steps")
	}
	if _, ok := steps["phase02_cli_session_manager.md"]; !ok {
		t.Fatal("missing phase02 step")
	}
}
