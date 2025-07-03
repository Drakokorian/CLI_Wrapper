package app

import (
	"fmt"

	"cli-wrapper/internal/logging"
	"cli-wrapper/internal/plugins"
)

// DetectCLITool returns the first available CLI tool.
func DetectCLITool() (string, error) {
	for _, p := range plugins.All() {
		if p.Detect() {
			return p.Name(), nil
		}
	}
	return "", fmt.Errorf("no supported CLI tool found")
}

// DetectCLITools returns all available CLI tools.
func DetectCLITools() ([]string, error) {
	tools := []string{}
	for _, p := range plugins.All() {
		if p.Detect() {
			tools = append(tools, p.Name())
		}
	}
	if len(tools) == 0 {
		return nil, fmt.Errorf("no supported CLI tool found")
	}
	return tools, nil
}

// InvokeTool runs the given CLI with args and logs the invocation.
func InvokeTool(tool string, args []string) error {
	logger, err := logging.New("info", "")
	if err != nil {
		return err
	}
	defer logger.Close()

	p, ok := plugins.Get(tool)
	if !ok {
		return fmt.Errorf("unknown tool %q", tool)
	}

	logger.Info(fmt.Sprintf("running %s %v", tool, args))
	if err := p.Invoke(args); err != nil {
		logger.Error(fmt.Sprintf("%s failed: %v", tool, err))
		return fmt.Errorf("run %s: %w", tool, err)
	}
	logger.Info(fmt.Sprintf("%s completed", tool))
	return nil
}
