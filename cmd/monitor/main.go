package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Carlos20052030/cloudmonitor/internal/checker"
	"github.com/Carlos20052030/cloudmonitor/internal/config"
	"github.com/Carlos20052030/cloudmonitor/internal/domain"
	"github.com/Carlos20052030/cloudmonitor/internal/storage"
)

// Store is the contract used by main. The concrete *storage.Store
// satisfies it implicitly — Go does not require "implements".
type Store interface {
	Save(ctx context.Context, r domain.Result) error
	Close() error
}

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Fprintln(os.Stderr, "config:", err)
		os.Exit(1)
	}

	store, err := storage.New("cloudmonitor.db")
	if err != nil {
		fmt.Fprintln(os.Stderr, "storage:", err)
		os.Exit(1)
	}
	defer store.Close()

	targets := make([]domain.Target, 0, len(cfg.Targets))
	for _, t := range cfg.Targets {
		targets = append(targets, domain.Target{Name: t.Name, URL: t.URL})
	}

	c := checker.New(time.Duration(cfg.RequestTimeoutSeconds) * time.Second)

	ctx := context.Background()
	results := c.CheckAll(ctx, targets)

	for _, r := range results {
		if err := store.Save(ctx, r); err != nil {
			fmt.Fprintln(os.Stderr, "save:", err)
			continue
		}
		fmt.Printf("%-25s %-5s HTTP %d %dms\n", r.TargetName, r.Status, r.StatusCode, r.LatencyMS)
	}
}
