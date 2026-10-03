package scheduler

import (
	"context"
	"time"

	"github.com/Carlos20052030/cloudmonitor/internal/checker"
	"github.com/Carlos20052030/cloudmonitor/internal/domain"
)

// Scheduler triggers checks at fixed intervals.
// The OnResult callback is injected by main — the scheduler does not
// know what happens to the result (save, notify, log).
type Scheduler struct {
	checker  *checker.Checker
	interval time.Duration
	targets  []domain.Target
	onResult func(domain.Result)
}

func New(c *checker.Checker, interval time.Duration, targets []domain.Target, onResult func(domain.Result)) *Scheduler {
	return &Scheduler{
		checker:  c,
		interval: interval,
		targets:  targets,
		onResult: onResult,
	}
}

// Run performs one immediate pass, then loops on the ticker.
// Returns cleanly when the context is cancelled.
func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.runOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.runOnce(ctx)
		}
	}
}

func (s *Scheduler) runOnce(ctx context.Context) {
	for _, r := range s.checker.CheckAll(ctx, s.targets) {
		s.onResult(r)
	}
}
