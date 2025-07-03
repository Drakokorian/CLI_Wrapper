package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"cli-wrapper/internal/app"
	"cli-wrapper/internal/history"
	"cli-wrapper/internal/logging"
)

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), `Usage:
  cliwrap [flags] run -model MODEL -prompt PROMPT
  cliwrap [flags] history export
  cliwrap [flags] history import

Flags:
`)
	flag.PrintDefaults()
}

func main() {
	conc := flag.Int("concurrency", 0, "set concurrency and update config")
	work := flag.String("workdir", "", "set working directory and update config")
	flag.Usage = usage
	flag.Parse()

	base, err := app.PrepareDirectories()
	if err != nil {
		log.Fatalf("prepare dirs: %v", err)
	}
	cfg, cfgErr := app.LoadConfig(base)
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	if cfg.LogPath == "" {
		cfg.LogPath = filepath.Join(base, "logs", "logs.txt")
	}

	if *conc > 0 {
		cfg.Concurrency = *conc
	}
	if *work != "" {
		cfg.WorkingDir = *work
	}
	if *conc > 0 || *work != "" {
		if err := app.SaveConfig(base, cfg); err != nil {
			log.Fatalf("save config: %v", err)
		}
	}

	logger, err := logging.New(cfg.LogLevel, cfg.LogPath)
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}
	defer logger.Close()
	if cfgErr != nil {
		logger.Error("load config: " + cfgErr.Error())
	}

	hist, err := history.New(base)
	if err != nil {
		logger.Error("history: " + err.Error())
		os.Exit(1)
	}
	defer hist.Close()

	mgr := app.NewSessionManager(logger, cfg.Concurrency, &cfg, hist)
	defer mgr.Close()

	args := flag.Args()
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	switch args[0] {
	case "run":
		runFS := flag.NewFlagSet("run", flag.ExitOnError)
		model := runFS.String("model", "", "model to run")
		prompt := runFS.String("prompt", "", "prompt text")
		runFS.Parse(args[1:])
		if *model == "" || *prompt == "" {
			fmt.Fprintln(os.Stderr, "model and prompt required")
			os.Exit(1)
		}
		if err := runPrompt(mgr, *model, *prompt); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	case "history":
		if len(args) < 2 {
			usage()
			os.Exit(1)
		}
		switch args[1] {
		case "export":
			data, err := hist.Export()
			if err != nil {
				logger.Error("export history: " + err.Error())
				os.Exit(1)
			}
			os.Stdout.Write(data)
		case "import":
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				logger.Error("read stdin: " + err.Error())
				os.Exit(1)
			}
			if err := hist.Import(data); err != nil {
				logger.Error("import history: " + err.Error())
				os.Exit(1)
			}
		default:
			usage()
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(1)
	}
}

func runPrompt(mgr *app.SessionManager, model, prompt string) error {
	id, err := mgr.AddSession(model, prompt, []string{prompt})
	if err != nil {
		return err
	}
	ch, err := mgr.OutputChannel(id)
	if err != nil {
		return err
	}
	for line := range ch {
		fmt.Println(line)
	}
	return nil
}
