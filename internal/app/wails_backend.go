package app

import (
	"context"
	"fmt"

	"cli-wrapper/internal/logging"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Backend exposes methods to the Wails frontend.
type Backend struct {
	ctx       context.Context
	mgr       *SessionManager
	logger    *logging.Logger
	cfg       *Config
	lastModel string
}

// NewBackend creates a Backend instance.
func NewBackend(mgr *SessionManager, logger *logging.Logger, cfg *Config) *Backend {
	return &Backend{mgr: mgr, logger: logger, cfg: cfg}
}

// Startup stores the context for runtime events.
func (b *Backend) Startup(ctx context.Context) {
	b.ctx = ctx
}

// RunPrompt executes the given model with the prompt.
func (b *Backend) RunPrompt(model, prompt string) (string, error) {
	if model == "" {
		return "", fmt.Errorf("model required")
	}
	if b.lastModel != "" && model != b.lastModel {
		if b.cfg.ModelAlerts {
			runtime.EventsEmit(b.ctx, "model:switched", map[string]string{"from": b.lastModel, "to": model})
		}
		_ = b.logger.ModelSwitch(b.lastModel, model, prompt)
	}
	b.lastModel = model
        id, err := b.mgr.AddSession(model, prompt, []string{prompt})
	if err != nil {
		return "", err
	}
	return id, nil
}
