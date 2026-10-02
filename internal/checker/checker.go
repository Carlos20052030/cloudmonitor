package checker

import (
	"net/http"
	"time"

	"github.com/Carlos20052030/cloudmonitor/internal/domain"
)

func Check(target domain.Target, timeout time.Duration) domain.Result {
	start := time.Now()
	result := domain.Result{
		TargetName: target.Name,
		URL:        target.URL,
		CheckedAt:  start,
	}

	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(target.URL)
	result.LatencyMS = time.Since(start).Milliseconds()

	if err != nil {
		result.Status = domain.StatusDOWN
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.Status = domain.StatusUP
	} else {
		result.Status = domain.StatusDOWN
	}

	return result
}