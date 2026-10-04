package scheduler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/Carlos20052030/cloudmonitor/internal/checker"
	"github.com/Carlos20052030/cloudmonitor/internal/domain"
)

func TestScheduler_RunsImmediately(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := checker.New(2 * time.Second)
	targets := []domain.Target{{Name: "a", URL: srv.URL}}

	var mu sync.Mutex
	var calls int

	s := New(c, 1*time.Hour, targets, func(r domain.Result) {
		mu.Lock()
		calls++
		mu.Unlock()
	})

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	go s.Run(ctx)
	<-ctx.Done()

	mu.Lock()
	defer mu.Unlock()
	if calls != 1 {
		t.Fatalf("expected exactly 1 call from immediate pass, got %d", calls)
	}
}

func TestScheduler_TickerFires(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := checker.New(2 * time.Second)
	targets := []domain.Target{{Name: "a", URL: srv.URL}}

	var mu sync.Mutex
	var calls int

	s := New(c, 50*time.Millisecond, targets, func(r domain.Result) {
		mu.Lock()
		calls++
		mu.Unlock()
	})

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()

	go s.Run(ctx)
	<-ctx.Done()

	mu.Lock()
	defer mu.Unlock()
	// 1 immediate + ~3 from ticker (50ms each within 180ms)
	if calls < 3 {
		t.Fatalf("expected at least 3 calls, got %d", calls)
	}
}

func TestScheduler_StopsOnContextCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := checker.New(2 * time.Second)
	targets := []domain.Target{{Name: "a", URL: srv.URL}}

	s := New(c, 10*time.Millisecond, targets, func(r domain.Result) {})

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		s.Run(ctx)
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case <-done:
		// ok, Run returned
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Run did not return after context cancel")
	}
}
