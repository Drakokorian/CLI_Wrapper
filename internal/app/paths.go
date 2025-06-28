package app

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

func appDirForOS(goos string) (string, error) {
	switch goos {
	case "windows":
		base := os.Getenv("APPDATA")
		if base == "" {
			return "", errors.New("APPDATA not set")
		}
		return filepath.Join(base, "ai-cli-ui"), nil
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, "Library", "Application Support", "ai-cli-ui"), nil
	default: // linux and others
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".config", "ai-cli-ui"), nil
	}
}

// AppDir returns the base directory for application data depending on the OS.
func AppDir() (string, error) {
	return appDirForOS(runtime.GOOS)
}

// PrepareDirectories ensures the logs, config, and state directories exist and
// returns the base path.
func PrepareDirectories() (string, error) {
	base, err := AppDir()
	if err != nil {
		return "", err
	}
	subdirs := []string{"logs", "config", "state"}
	for _, d := range subdirs {
		p := filepath.Join(base, d)
		if err := os.MkdirAll(p, 0o755); err != nil {
			return "", err
		}
	}
	return base, nil
}
