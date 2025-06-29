package main

import (
	"log"
	"path/filepath"

	"github.com/wailsapp/wails/v2"

	"cli-wrapper/internal/app"
	"cli-wrapper/internal/history"
	"cli-wrapper/internal/logging"
)

func main() {
	base, err := app.AppDir()
	if err != nil {
		log.Printf("app dir: %v", err)
		return
	}
	logger, err := logging.NewWithPath(filepath.Join(base, "logs", "logs.txt"))
	if err != nil {
		log.Printf("init logger: %v", err)
		return
	}
	defer logger.Close()

	if err := app.PrepareDirectories(); err != nil {
		logger.Error("prepare dirs: " + err.Error())
		return
	}

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
	if err := wails.Run(&wails.Options{Bind: []interface{}{backend}, OnStartup: backend.Startup}); err != nil {
		logger.Error(err.Error())
	}
}
