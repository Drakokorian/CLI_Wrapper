package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds persistent application settings.
// Config struct includes both runtime and UI preferences.
// Config holds persistent application settings including resource thresholds.
type Config struct {
	Concurrency     int     `json:"concurrency"`
	Theme           string  `json:"theme"`
	ModelAlerts     bool    `json:"modelAlerts"`
	CPUThreshold    float64 `json:"cpuThreshold"`
	MemoryThreshold float64 `json:"memoryThreshold"`
	PollInterval    int     `json:"pollInterval"`
	WorkingDir      string  `json:"workingDir"`
	LogLevel        string  `json:"logLevel"`
	LogPath         string  `json:"logPath"`
}

// ValidTheme reports whether the provided theme value is supported.
func ValidTheme(t string) bool {
	return t == "light" || t == "dark"
}

const defaultTheme = "light"
const defaultCPUThreshold = 35.0
const defaultMemoryThreshold = 35.0
const defaultPollInterval = 2
const defaultLogLevel = "info"

func defaultLogPath(baseDir string) string {
	return filepath.Join(baseDir, "logs", "logs.txt")
}

func defaultWorkingDir() string {
	if dir, err := os.UserHomeDir(); err == nil {
		return dir
	}
	return "."
}

// LoadConfig reads configuration from the config directory.
func LoadConfig(baseDir string) (Config, error) {
	path := filepath.Join(baseDir, "config", "config.json")
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{
				Concurrency:     1,
				Theme:           defaultTheme,
				ModelAlerts:     true,
				CPUThreshold:    defaultCPUThreshold,
				MemoryThreshold: defaultMemoryThreshold,
				PollInterval:    defaultPollInterval,
				WorkingDir:      defaultWorkingDir(),
				LogLevel:        defaultLogLevel,
				LogPath:         defaultLogPath(baseDir),
			}, nil
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
	if cfg.CPUThreshold <= 0 || cfg.CPUThreshold > 100 {
		cfg.CPUThreshold = defaultCPUThreshold
	}
	if cfg.MemoryThreshold <= 0 || cfg.MemoryThreshold > 100 {
		cfg.MemoryThreshold = defaultMemoryThreshold
	}
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = defaultPollInterval
	}
	if !cfg.ModelAlerts {
		cfg.ModelAlerts = true
	}
	switch cfg.LogLevel {
	case "error", "info", "debug":
	default:
		cfg.LogLevel = defaultLogLevel
	}
	if cfg.LogPath == "" {
		cfg.LogPath = defaultLogPath(baseDir)
	}
	if cfg.WorkingDir == "" {
		cfg.WorkingDir = defaultWorkingDir()
	} else if info, err := os.Stat(cfg.WorkingDir); err != nil || !info.IsDir() {
		cfg.WorkingDir = defaultWorkingDir()
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
	if cfg.CPUThreshold <= 0 || cfg.CPUThreshold > 100 {
		return fmt.Errorf("invalid CPU threshold %f", cfg.CPUThreshold)
	}
	if cfg.MemoryThreshold <= 0 || cfg.MemoryThreshold > 100 {
		return fmt.Errorf("invalid memory threshold %f", cfg.MemoryThreshold)
	}
	if cfg.PollInterval <= 0 {
		return fmt.Errorf("invalid poll interval %d", cfg.PollInterval)
	}
	if !cfg.ModelAlerts {
		// allow both true and false; no validation necessary
	}
	switch cfg.LogLevel {
	case "error", "info", "debug":
	default:
		return fmt.Errorf("invalid log level %q", cfg.LogLevel)
	}
	if cfg.LogPath == "" {
		return fmt.Errorf("log path required")
	}
	if cfg.WorkingDir == "" {
		return fmt.Errorf("working directory required")
	}
	if info, err := os.Stat(cfg.WorkingDir); err != nil || !info.IsDir() {
		return fmt.Errorf("invalid working directory %q", cfg.WorkingDir)
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
