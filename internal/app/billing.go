package app

import (
	"fmt"

	"cli-wrapper/internal/plugins"
)

// BillingUsage holds credit usage information returned by the CLI.
// BillingUsage is an alias for plugins.BillingUsage for backward compatibility.
type BillingUsage = plugins.BillingUsage

// BillingURL returns the billing page URL for the given CLI tool.
func BillingURL(tool string) (string, error) {
	p, ok := plugins.Get(tool)
	if !ok {
		return "", fmt.Errorf("unknown tool %q", tool)
	}
	return p.BillingURL(), nil
}

// FetchUsage executes the CLI to retrieve usage information if supported.
// It returns an error if the CLI does not provide usage details.
func FetchUsage(tool string) (*BillingUsage, error) {
	p, ok := plugins.Get(tool)
	if !ok {
		return nil, fmt.Errorf("unsupported tool %q", tool)
	}
	return p.Usage()
}
