package plugins

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type openAIPlugin struct{}

func (openAIPlugin) Name() string { return "openai" }

func (openAIPlugin) Detect() bool {
	_, err := exec.LookPath("openai")
	return err == nil
}

func (openAIPlugin) Invoke(args []string) error {
	cmd := exec.Command("openai", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("openai: %w", err)
	}
	fmt.Print(string(out))
	return nil
}

func (openAIPlugin) BillingURL() string {
	return "https://platform.openai.com/account/billing"
}

func (openAIPlugin) Usage() (*BillingUsage, error) {
	out, err := exec.Command("openai", "api", "request", "GET", "/dashboard/billing/credit_grants").Output()
	if err != nil {
		return nil, fmt.Errorf("openai usage: %w", err)
	}
	var u BillingUsage
	if err := json.Unmarshal(out, &u); err != nil {
		return nil, fmt.Errorf("parse usage: %w", err)
	}
	return &u, nil
}

func init() { Register(openAIPlugin{}) }
