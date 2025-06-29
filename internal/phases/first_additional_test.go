package phases

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFirstStepErrors(t *testing.T) {
	_, err := FirstStep(filepath.Join(t.TempDir(), "missing.md"))
	if err == nil {
		t.Fatal("expected error")
	}

	path := filepath.Join(t.TempDir(), "p.md")
	if err := os.WriteFile(path, []byte("## Detailed Steps\n## Next"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if _, err := FirstStep(path); err == nil {
		t.Fatal("expected no bullet error")
	}
}

func TestPhaseFirstStepsMixed(t *testing.T) {
	dir := t.TempDir()
	good := filepath.Join(dir, "phase01.md")
	bad := filepath.Join(dir, "phase99.md")
	os.WriteFile(good, []byte("## Detailed Steps\n- step1"), 0o644)
	os.WriteFile(bad, []byte("## Detailed Steps\nnot-a-bullet"), 0o644)
	_, err := PhaseFirstSteps(dir)
	if err == nil {
		t.Fatal("expected error from bad file")
	}
}
