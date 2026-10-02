package checker

import (
	"context"
	"net/http"
	"time"

	"github.com/Carlos20052030/cloudmonitor/internal/domain"
)

type Checker struct {
	client *http.Client
}

func New(timeout time.Duration) *Checker {
	return &Checker{
		client: &http.Client{Timeout: timeout},
	}
}

func (c *Checker) Check(ctx context.Context, target domain.Target) domain.Result {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target.URL, nil)
	if err != nil {
		return domain.Result{
			TargetName: target.Name,
			URL:        target.URL,
			Status:     domain.StatusDOWN,
			Error:      err.Error(),
			CheckedAt:  time.Now(),
			LatencyMS:  time.Since(start).Milliseconds(),
		}
	}

	resp, err := c.client.Do(req)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return domain.Result{
			TargetName: target.Name,
			URL:        target.URL,
			Status:     domain.StatusDOWN,
			Error:      err.Error(),
			CheckedAt:  time.Now(),
			LatencyMS:  latency,
		}
	}
	defer resp.Body.Close()

	status := domain.StatusDOWN
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		status = domain.StatusUP
	}

	return domain.Result{
		TargetName: target.Name,
		URL:        target.URL,
		Status:     status,
		StatusCode: resp.StatusCode,
		LatencyMS:  latency,
		CheckedAt:  time.Now(),
	}
}
