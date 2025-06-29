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
	base, err := app.PrepareDirectories()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := app.LoadConfig(base)
	if err != nil {
		log.Fatal(err)
	}
        logger, err := logging.NewWithPath(filepath.Join(base, "logs", "logs.txt"))
        if err != nil {
                log.Fatal(err)
        }
        defer logger.Close()
        hist, err := history.New(base)
        if err != nil {
                log.Fatal(err)
        }
        defer hist.Close()

        mgr := app.NewSessionManager(base, logger, cfg.Concurrency, &cfg, hist)
        defer mgr.Close()
        backend := app.NewBackend(mgr, logger, &cfg)
	err = wails.Run(&wails.Options{Bind: []interface{}{backend}, OnStartup: backend.Startup})
	if err != nil {
		log.Fatal(err)
	}
}
