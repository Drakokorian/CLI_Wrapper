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
	base, err := app.PrepareDirectories()
	if err != nil {
		log.Printf("prepare dirs: %v", err)
		return
	}
	logger, err := logging.NewWithPath(filepath.Join(base, "logs", "logs.txt"))
	if err != nil {
		log.Printf("init logger: %v", err)
		return
	}
	defer logger.Close()

	cfg, err := app.LoadConfig(base)
	if err != nil {
		logger.Error("load config: " + err.Error())
		return
	}
	hist, err := history.New(base)
	if err != nil {
		logger.Error("history: " + err.Error())
		return
	}
	defer hist.Close()

	mgr := app.NewSessionManager(base, logger, cfg.Concurrency, &cfg, hist)
	defer mgr.Close()
	backend := app.NewBackend(mgr, logger, &cfg)
	opts := &options.App{Bind: []interface{}{backend}, OnStartup: backend.Startup}
	if err := wails.Run(opts); err != nil {
		logger.Error(err.Error())
	}
}
