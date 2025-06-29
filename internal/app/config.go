package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds persistent application settings.
// Config struct includes both runtime and UI preferences.
type Config struct {
	Concurrency int    `json:"concurrency"`
	Theme       string `json:"theme"`
	ModelAlerts bool   `json:"modelAlerts"`
}

// ValidTheme reports whether the provided theme value is supported.
func ValidTheme(t string) bool {
	return t == "light" || t == "dark"
}

const defaultTheme = "light"

// LoadConfig reads configuration from the config directory.
func LoadConfig(baseDir string) (Config, error) {
	path := filepath.Join(baseDir, "config", "config.json")
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{Concurrency: 1, Theme: defaultTheme, ModelAlerts: true}, nil
		}
		return Config{}, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()
	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("decode config: %w", err)
	}
	if cfg.Concurrency < 1 || cfg.Concurrency > 5 {
		cfg.Concurrency = 1
	}
	if !ValidTheme(cfg.Theme) {
		cfg.Theme = defaultTheme
	}
	if !cfg.ModelAlerts {
		// missing field defaults to true when false and file existed
		cfg.ModelAlerts = true
	}
	return cfg, nil
}

// SaveConfig writes configuration to the config directory.
func SaveConfig(baseDir string, cfg Config) error {
	if cfg.Concurrency < 1 || cfg.Concurrency > 5 {
		return fmt.Errorf("invalid concurrency %d", cfg.Concurrency)
	}
	if !ValidTheme(cfg.Theme) {
		return fmt.Errorf("invalid theme %q", cfg.Theme)
	}
	if !cfg.ModelAlerts {
		// allow both true and false; no validation necessary
	}
	path := filepath.Join(baseDir, "config", "config.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("open config for write: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cfg); err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	return nil
}
