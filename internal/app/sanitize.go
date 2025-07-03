package app

import (
	"fmt"
	"regexp"
	"strings"
)

// modelPattern defines allowable characters for model identifiers.
// Models may contain letters, numbers, dashes, underscores and periods.
var modelPattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

// SanitizeModel validates allowed characters in model names.
func SanitizeModel(model string) (string, error) {
	if !modelPattern.MatchString(model) {
		return "", fmt.Errorf("invalid model name")
	}
	return model, nil
}

// SanitizePrompt removes control characters and leading dashes to avoid argument injection.
func SanitizePrompt(p string) string {
	p = strings.ReplaceAll(p, "\n", " ")
	p = strings.ReplaceAll(p, "\r", " ")
	p = strings.TrimSpace(p)
	for strings.HasPrefix(p, "-") {
		p = strings.TrimPrefix(p, "-")
		p = "_" + p
	}
	return p
}
