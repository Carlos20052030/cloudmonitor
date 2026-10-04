package storage

import (
	"context"
	"testing"
	"time"

	"github.com/Carlos20052030/cloudmonitor/internal/domain"
)

func TestSave(t *testing.T) {
	s, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to open store: %v", err)
	}
	defer func() { _ = s.Close() }()

	r := domain.Result{
		TargetName: "example",
		URL:        "https://example.com",
		Status:     domain.StatusUP,
		StatusCode: 200,
		LatencyMS:  42,
		CheckedAt:  time.Now(),
	}

	if err := s.Save(context.Background(), r); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	var count int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM checks").Scan(&count); err != nil {
		t.Fatalf("failed to query: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 row, got %d", count)
	}
}

func TestSave_DOWN(t *testing.T) {
	s, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to open store: %v", err)
	}
	defer func() { _ = s.Close() }()

	r := domain.Result{
		TargetName: "broken",
		URL:        "https://broken.invalid",
		Status:     domain.StatusDOWN,
		StatusCode: 0,
		Error:      "connection refused",
		CheckedAt:  time.Now(),
	}

	if err := s.Save(context.Background(), r); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	var status, errMsg string
	if err := s.db.QueryRow("SELECT status, error FROM checks LIMIT 1").Scan(&status, &errMsg); err != nil {
		t.Fatalf("failed to query: %v", err)
	}
	if status != "DOWN" {
		t.Fatalf("expected DOWN, got %s", status)
	}
	if errMsg != "connection refused" {
		t.Fatalf("expected error preserved, got %q", errMsg)
	}
}

func TestRecentChecks(t *testing.T) {
	s, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to open store: %v", err)
	}
	defer func() { _ = s.Close() }()

	now := time.Now()
	for i := 0; i < 5; i++ {
		r := domain.Result{
			TargetName: "example",
			URL:        "https://example.com",
			Status:     domain.StatusUP,
			StatusCode: 200,
			LatencyMS:  int64(i * 10),
			CheckedAt:  now.Add(time.Duration(i) * time.Second),
		}
		if err := s.Save(context.Background(), r); err != nil {
			t.Fatalf("failed to save: %v", err)
		}
	}

	got, err := s.RecentChecks(context.Background(), "example", 3)
	if err != nil {
		t.Fatalf("failed to query: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 records, got %d", len(got))
	}
	// Most recent first: latencies should be 40, 30, 20
	if got[0].LatencyMS != 40 {
		t.Fatalf("expected first latency 40, got %d", got[0].LatencyMS)
	}
	if got[2].LatencyMS != 20 {
		t.Fatalf("expected third latency 20, got %d", got[2].LatencyMS)
	}
}
