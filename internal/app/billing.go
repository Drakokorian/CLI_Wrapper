package app

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// BillingUsage holds credit usage information returned by the CLI.
type BillingUsage struct {
	TotalGranted   float64 `json:"total_granted"`
	TotalUsed      float64 `json:"total_used"`
	TotalAvailable float64 `json:"total_available"`
}

// BillingURL returns the billing page URL for the given CLI tool.
func BillingURL(tool string) (string, error) {
	switch tool {
	case "openai":
		return "https://platform.openai.com/account/billing", nil
	case "gemini":
		return "https://makersuite.google.com/app/apikey", nil
	default:
		return "", fmt.Errorf("unknown tool %q", tool)
	}
}

// FetchUsage executes the CLI to retrieve usage information if supported.
// It returns an error if the CLI does not provide usage details.
func FetchUsage(tool string) (*BillingUsage, error) {
	switch tool {
	case "openai":
		out, err := exec.Command("openai", "api", "request", "GET", "/dashboard/billing/credit_grants").Output()
		if err != nil {
			return nil, fmt.Errorf("openai usage: %w", err)
		}
		var res BillingUsage
		if err := json.Unmarshal(out, &res); err != nil {
			return nil, fmt.Errorf("parse openai usage: %w", err)
		}
		return &res, nil
	case "gemini":
		out, err := exec.Command("gemini", "usage").Output()
		if err != nil {
			return nil, fmt.Errorf("gemini usage: %w", err)
		}
		var res BillingUsage
		if err := json.Unmarshal(out, &res); err != nil {
			return nil, fmt.Errorf("parse gemini usage: %w", err)
		}
		return &res, nil
	default:
		return nil, fmt.Errorf("unsupported tool %q", tool)
	}
}
