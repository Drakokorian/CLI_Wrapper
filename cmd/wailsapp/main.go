package main

import (
	"log"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"

	"cli-wrapper/internal/app"
	"cli-wrapper/internal/history"
	"cli-wrapper/internal/logging"
)

func main() {
	// base is the root path containing config, logs and state directories.
	base, err := app.PrepareDirectories()
	if err != nil {
		log.Printf("prepare dirs: %v", err)
		return
	}
	cfg, cfgErr := app.LoadConfig(base)
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	if cfg.LogPath == "" {
		cfg.LogPath = filepath.Join(base, "logs", "logs.txt")
	}
	logger, err := logging.New(cfg.LogLevel, cfg.LogPath)
	if err != nil {
		log.Printf("init logger: %v", err)
		return
	}
	defer logger.Close()
	if cfgErr != nil {
		logger.Error("load config: " + cfgErr.Error())
	}
	hist, err := history.New(base)
	if err != nil {
		logger.Error("history: " + err.Error())
		return
	}
	defer hist.Close()

	mgr := app.NewSessionManager(logger, cfg.Concurrency, &cfg, hist)
	defer mgr.Close()
	backend := app.NewBackend(mgr, logger, &cfg)
	opts := &options.App{Bind: []interface{}{backend}, OnStartup: backend.Startup}
	if err := wails.Run(opts); err != nil {
		logger.Error(err.Error())
	}
}
