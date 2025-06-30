package plugins

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type geminiPlugin struct{}

func (geminiPlugin) Name() string { return "gemini" }

func (geminiPlugin) Detect() bool {
	_, err := exec.LookPath("gemini")
	return err == nil
}

func (geminiPlugin) Invoke(args []string) error {
	cmd := exec.Command("gemini", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gemini: %w", err)
	}
	fmt.Print(string(out))
	return nil
}

func (geminiPlugin) BillingURL() string {
	return "https://makersuite.google.com/app/apikey"
}

func (geminiPlugin) Usage() (*BillingUsage, error) {
	out, err := exec.Command("gemini", "usage").Output()
	if err != nil {
		return nil, fmt.Errorf("gemini usage: %w", err)
	}
	var u BillingUsage
	if err := json.Unmarshal(out, &u); err != nil {
		return nil, fmt.Errorf("parse usage: %w", err)
	}
	return &u, nil
}

func init() { Register(geminiPlugin{}) }
