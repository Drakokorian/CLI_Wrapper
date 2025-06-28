package phases

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FirstStep returns the first bullet point from the Detailed Steps section
// of the specified phase file.
func FirstStep(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open phase file: %w", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	inSection := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "## ") && inSection {
			break
		}
		if strings.EqualFold(strings.TrimSpace(line), "## Detailed Steps") {
			inSection = true
			continue
		}
		if inSection && strings.HasPrefix(strings.TrimSpace(line), "-") {
			return strings.TrimSpace(strings.TrimPrefix(line, "-")), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("scan phase file: %w", err)
	}
	return "", fmt.Errorf("no bullet found in %s", path)
}

// PhaseFirstSteps reads all phase markdown files in the given directory
// and returns a map of filename to first step bullet.
func PhaseFirstSteps(dir string) (map[string]string, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "phase*.md"))
	if err != nil {
		return nil, fmt.Errorf("glob phases: %w", err)
	}
	steps := make(map[string]string)
	for _, p := range paths {
		step, err := FirstStep(p)
		if err != nil {
			return nil, err
		}
		steps[filepath.Base(p)] = step
	}
	return steps, nil
}
