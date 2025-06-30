package main

import (
	"log"
	"path/filepath"

	"cli-wrapper/internal/app"
	"cli-wrapper/internal/history"
	"cli-wrapper/internal/logging"
	"cli-wrapper/internal/server"
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

	mgr := app.NewSessionManager(base, logger, cfg.Concurrency, &cfg, hist)
	defer mgr.Close()

	srv := server.New(mgr, logger, base, &cfg, hist)
	if err := srv.Start(":8080"); err != nil {
		logger.Error(err.Error())
	}
}
