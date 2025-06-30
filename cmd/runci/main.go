package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"os/exec"

	"cli-wrapper/internal/logging"
)

func run(logger *logging.Logger, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	buf := bytes.TrimSpace(out)
	if len(buf) > 0 {
		scanner := bufio.NewScanner(bytes.NewReader(buf))
		for scanner.Scan() {
			logger.Info(name + ": " + scanner.Text())
		}
	}
	if err != nil {
		logger.Error(name + " failed: " + err.Error())
	} else {
		logger.Info(name + " succeeded")
	}
	return err
}

func main() {
	logger, err := logging.New("info", "")
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}
	defer logger.Close()

	if err := run(logger, "go", "vet", "./..."); err != nil {
		os.Exit(1)
	}
	if err := run(logger, "go", "test", "./..."); err != nil {
		os.Exit(1)
	}
	if err := run(logger, "go", "build", "./..."); err != nil {
		os.Exit(1)
	}
}
