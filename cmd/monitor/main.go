package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Carlos20052030/cloudmonitor/internal/checker"
	"github.com/Carlos20052030/cloudmonitor/internal/config"
	"github.com/Carlos20052030/cloudmonitor/internal/domain"
	"github.com/Carlos20052030/cloudmonitor/internal/scheduler"
	"github.com/Carlos20052030/cloudmonitor/internal/storage"
)

// Store is the contract used by main. The concrete *storage.Store
// satisfies it implicitly — Go does not require "implements".
type Store interface {
	Save(ctx context.Context, r domain.Result) error
	Close() error
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	cfg, err := config.Load("config.yaml")
	if err != nil {
		slog.Error("config", "err", err)
		os.Exit(1)
	}

	store, err := storage.New("cloudmonitor.db")
	if err != nil {
		slog.Error("storage", "err", err)
		os.Exit(1)
	}
	defer store.Close()

	targets := make([]domain.Target, 0, len(cfg.Targets))
	for _, t := range cfg.Targets {
		targets = append(targets, domain.Target{Name: t.Name, URL: t.URL})
	}

	c := checker.New(time.Duration(cfg.RequestTimeoutSeconds) * time.Second)

	onResult := func(r domain.Result) {
		slog.Info("check",
			"target", r.TargetName,
			"status", r.Status,
			"status_code", r.StatusCode,
			"latency_ms", r.LatencyMS,
		)
		if err := store.Save(context.Background(), r); err != nil {
			slog.Error("save", "err", err)
		}
	}

	interval := time.Duration(cfg.CheckIntervalSeconds) * time.Second
	sched := scheduler.New(c, interval, targets, onResult)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	slog.Info("monitor starting", "targets", len(targets), "interval", interval)
	sched.Run(ctx)
	slog.Info("monitor stopped")
}
