package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Carlos20052030/cloudmonitor/internal/checker"
	"github.com/Carlos20052030/cloudmonitor/internal/config"
	"github.com/Carlos20052030/cloudmonitor/internal/domain"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	timeout := time.Duration(cfg.RequestTimeoutSeconds) * time.Second

	for _, target := range cfg.Targets {
        result := checker.Check(domain.Target{Name: target.Name, URL: target.URL}, timeout)
		fmt.Printf("URL: %s\n", result.URL)
		fmt.Printf("Status: %s\n", result.Status)
		if result.Error != "" {
			fmt.Printf("Error: %s\n", result.Error)
		} else {
			fmt.Printf("HTTP: %d\n", result.StatusCode)
		}
		fmt.Printf("Latency: %dms\n", result.LatencyMS)
		fmt.Println("---")
	}
}