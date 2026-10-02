package main

import (
	"fmt"
	"os"

	"github.com/Carlos20052030/cloudmonitor/internal/config"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	fmt.Printf("interval: %ds\n", cfg.CheckIntervalSeconds)
	fmt.Printf("timeout: %ds\n", cfg.RequestTimeoutSeconds)
	fmt.Printf("targets: %d\n", len(cfg.Targets))
	for _, t := range cfg.Targets {
		fmt.Printf("  - %s: %s\n", t.Name, t.URL)
	}
}