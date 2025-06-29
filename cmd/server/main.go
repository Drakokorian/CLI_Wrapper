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

        srv := server.New(mgr, logger, base, &cfg, hist)
	if err := srv.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
