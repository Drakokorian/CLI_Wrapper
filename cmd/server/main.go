package main

import (
	"log"

	"cli-wrapper/internal/app"
	"cli-wrapper/internal/history"
	"cli-wrapper/internal/logging"
	"cli-wrapper/internal/server"
)

func main() {
	base, err := app.PrepareDirectories()
	if err != nil {
		log.Printf("prepare dirs: %v", err)
		return
	}
	logger, err := logging.New()
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

	srv := server.New(mgr, logger, base, &cfg, hist)
	if err := srv.Start(":8080"); err != nil {
		logger.Error(err.Error())
	}
}
