package app

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"cli-wrapper/internal/logging"
)

// DetectCLITool returns the first available CLI tool.
func DetectCLITool() (string, error) {
	if _, err := exec.LookPath("openai"); err == nil {
		return "openai", nil
	}
	if _, err := exec.LookPath("gemini"); err == nil {
		return "gemini", nil
	}
	return "", fmt.Errorf("no supported CLI tool found")
}

// DetectCLITools returns all available CLI tools.
func DetectCLITools() ([]string, error) {
	tools := []string{}
	if _, err := exec.LookPath("openai"); err == nil {
		tools = append(tools, "openai")
	}
	if _, err := exec.LookPath("gemini"); err == nil {
		tools = append(tools, "gemini")
	}
	if len(tools) == 0 {
		return nil, fmt.Errorf("no supported CLI tool found")
	}
	return tools, nil
}

// InvokeTool runs the given CLI with args and logs the invocation.
func InvokeTool(tool string, args []string, baseDir string) error {
	logPath := filepath.Join(baseDir, "logs", "logs.txt")
	logger, err := logging.NewWithPath(logPath)
	if err != nil {
		return err
	}
	defer logger.Close()

	cmd := exec.Command(tool, args...)
	logger.Info(fmt.Sprintf("running %s %v", tool, args))
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error(fmt.Sprintf("%s failed: %v", tool, err))
		return fmt.Errorf("run %s: %w", tool, err)
	}
	logger.Info(fmt.Sprintf("%s output: %s", tool, string(out)))
	return nil
}
