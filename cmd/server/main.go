package main

import (
	"log"
	"path/filepath"

	"cli-wrapper/internal/app"
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

	mgr := app.NewSessionManager(base, logger, cfg.Concurrency, &cfg)
	defer mgr.Close()

	srv := server.New(mgr, logger, base, &cfg)
	if err := srv.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
