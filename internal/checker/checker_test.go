package checker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Carlos20052030/cloudmonitor/internal/domain"
)

// TestCheck_UP covers the happy path: the server responds 200.
func TestCheck_UP(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(2 * time.Second)
	r := c.Check(context.Background(), domain.Target{Name: "t", URL: srv.URL})

	if r.Status != domain.StatusUP {
		t.Fatalf("expected UP, got %s", r.Status)
	}
	if r.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", r.StatusCode)
	}
	if r.Error != "" {
		t.Fatalf("expected no error, got %q", r.Error)
	}
}

// TestCheck_DOWN covers a server that responds 500.
func TestCheck_DOWN(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := New(2 * time.Second)
	r := c.Check(context.Background(), domain.Target{Name: "t", URL: srv.URL})

	if r.Status != domain.StatusDOWN {
		t.Fatalf("expected DOWN, got %s", r.Status)
	}
	if r.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", r.StatusCode)
	}
}

// TestCheck_Timeout covers a server slower than the client timeout.
func TestCheck_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(100 * time.Millisecond)
	r := c.Check(context.Background(), domain.Target{Name: "t", URL: srv.URL})

	if r.Status != domain.StatusDOWN {
		t.Fatalf("expected DOWN on timeout, got %s", r.Status)
	}
	if r.Error == "" {
		t.Fatal("expected timeout error, got empty")
	}
}

// TestCheck_BadURL covers an invalid URL (request build error).
func TestCheck_BadURL(t *testing.T) {
	c := New(2 * time.Second)
	r := c.Check(context.Background(), domain.Target{Name: "t", URL: "ht!tp://invalido"})

	if r.Status != domain.StatusDOWN {
		t.Fatalf("expected DOWN, got %s", r.Status)
	}
	if r.Error == "" {
		t.Fatal("expected error, got empty")
	}
}

// TestCheckAll covers parallel dispatch across multiple targets.
func TestCheckAll(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(2 * time.Second)
	targets := []domain.Target{
		{Name: "a", URL: srv.URL},
		{Name: "b", URL: srv.URL},
		{Name: "c", URL: srv.URL},
	}

	results := c.CheckAll(context.Background(), targets)

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Status != domain.StatusUP {
			t.Fatalf("expected UP for %s, got %s", r.TargetName, r.Status)
		}
	}
}
