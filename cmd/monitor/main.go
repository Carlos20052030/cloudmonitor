package main

import (
	"context"
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

	targets := make([]domain.Target, 0, len(cfg.Targets))
	for _, t := range cfg.Targets {
		targets = append(targets, domain.Target{Name: t.Name, URL: t.URL})
	}

	c := checker.New(time.Duration(cfg.RequestTimeoutSeconds) * time.Second)

	start := time.Now()
	results := c.CheckAll(context.Background(), targets)
	elapsed := time.Since(start)

	for _, r := range results {
		fmt.Printf("%-25s %-5s HTTP %d %dms\n", r.TargetName, r.Status, r.StatusCode, r.LatencyMS)
	}
	fmt.Printf("\ntotal: %dms (parallel)\n", elapsed.Milliseconds())
}
