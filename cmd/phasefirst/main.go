package main

import (
	"encoding/json"
	"fmt"
	"log"

	"cli-wrapper/internal/logging"
	"cli-wrapper/internal/phases"
)

func main() {
	logger, err := logging.New()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Close()

	steps, err := phases.PhaseFirstSteps("docs/phases")
	if err != nil {
		logger.Error(fmt.Sprintf("collect steps: %v", err))
		log.Fatal(err)
	}
	data, err := json.MarshalIndent(steps, "", "  ")
	if err != nil {
		logger.Error(fmt.Sprintf("marshal output: %v", err))
		log.Fatal(err)
	}
	logger.Info("extracted first steps successfully")
	fmt.Println(string(data))
}
