# Architecture — Uptime Monitor V1

## 1. Overview

Uptime Monitor is a long-running service that periodically checks a set of
HTTP targets, records the outcome of each check, and notifies a Telegram
chat when a target changes state (UP → DOWN or DOWN → UP).

The system is intentionally small, dependency-light, and built on the Go
standard library wherever possible. It is designed to run as a single
process on a small VPS under systemd.

```text
                 cmd/monitor
                      │
                      │ wires dependencies
                      ▼
               ┌───────────────┐
               │   scheduler   │
               └───────┬───────┘
                       │ every N seconds
                       ▼
               ┌───────────────┐
               │    checker    │  (parallel, goroutines)
               └───────┬───────┘
                       │ produces domain.Result
                       ▼
         ┌─────────────┴─────────────┐
         ▼                           ▼
 ┌───────────────┐           ┌───────────────┐
 │    storage    │           │   notifier    │
 │   (SQLite)    │           │  (Telegram)   │
 └───────────────┘           └───────────────┘
````

## 2. Design Principles

1. **Single responsibility per package.** Each `internal/` package does one
   thing and exposes a small surface.

2. **Dependencies point inward.** All packages may depend on
   `internal/domain`; `domain` depends on nothing. `checker` does not import
   `storage`. `storage` does not import `checker`.

3. **Interfaces at the boundary.** `storage` and `notifier` are consumed
   through interfaces defined by the consumer, not by the implementation.
   This keeps the core testable and swappable.

4. **Standard library first.** `net/http`, `database/sql`, `log/slog`,
   `context`, `time`, and `os/signal` are preferred. Third-party code is used
   only where the standard library cannot reasonably do the job.

5. **Fail loud, exit clean.** Errors during startup cause a non-zero exit
   with a structured log line. Signals cause a graceful shutdown through
   `context`.

6. **No global state.** Every dependency is injected by `cmd/monitor`.

## 3. Component Responsibilities

| Package              | Responsibility                                                     | Depends on                                     |
| -------------------- | ------------------------------------------------------------------ | ---------------------------------------------- |
| `cmd/monitor`        | Read config, wire dependencies, start scheduler, handle signals    | all internal packages                          |
| `internal/domain`    | Shared types: `Target`, `Result`, `Status`                         | nothing                                        |
| `internal/config`    | Load and validate `config.yaml`                                    | `gopkg.in/yaml.v3`, `domain`                   |
| `internal/checker`   | Perform HTTP checks in parallel, measure latency, produce `Result` | `net/http`, `context`, `domain`                |
| `internal/storage`   | Persist `Result` rows to SQLite, run migrations                    | `database/sql`, `modernc.org/sqlite`, `domain` |
| `internal/notifier`  | Send Telegram messages on state transitions                        | Telegram Bot API, `domain`                     |
| `internal/scheduler` | Trigger checks on a fixed interval, forward results                | `context`, `time`, `checker`, `domain`         |

## 4. Domain Model

```go
package domain

import "time"

type Status string

const (
    StatusUP   Status = "UP"
    StatusDOWN Status = "DOWN"
)

type Target struct {
    Name string
    URL  string
}

type Result struct {
    TargetName string
    URL        string
    Status     Status
    StatusCode int
    LatencyMS  int64
    Error      string
    CheckedAt  time.Time
}
```

All packages use these types. No package redefines them.

## 5. Interfaces

```go
// Consumed by cmd/monitor and scheduler.
type Store interface {
    Save(ctx context.Context, r domain.Result) error
    Close() error
}

// Consumed by cmd/monitor when forwarding results.
type Notifier interface {
    Notify(ctx context.Context, r domain.Result) error
}
```

Interfaces are declared where they are used, not where they are implemented.
This is the Go idiom.

## 6. Execution Flow

1. `cmd/monitor/main.go` starts.
2. `config.Load("config.yaml")` reads and parses the YAML file into
   `config.Config`.
3. `storage.New("uptime.db")` opens the SQLite database and runs migrations.
4. If `config.Telegram.Enabled` is true, `notifier.NewTelegram(...)` is
   created.
5. `checker.New(timeout)` is created with the configured request timeout.
6. `scheduler.New(checker, interval, targets, onResult)` is created.
   `onResult` is a closure defined in `main` that:

   * logs the result with `slog`;
   * persists it via `storage.Save`;
   * forwards it to the notifier if enabled.
7. `signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)`
   creates a context that is cancelled on shutdown.
8. `scheduler.Run(ctx)`:

   * runs one immediate check pass;
   * then loops on a `time.Ticker`;
   * exits when `ctx.Done()` fires.
9. On shutdown, `storage.Close()` is deferred and runs before process exit.

## 7. Storage Schema

SQLite file: `uptime.db`

```sql
CREATE TABLE IF NOT EXISTS checks (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    target_name  TEXT     NOT NULL,
    url          TEXT     NOT NULL,
    status       TEXT     NOT NULL,
    status_code  INTEGER,
    latency_ms   INTEGER,
    error        TEXT,
    checked_at   DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_checks_target_time
    ON checks (target_name, checked_at);
```

The index supports the most common query pattern: "last N checks for a
given target".

Migrations live in `migrations/001_init.sql` and are executed idempotently
at startup.

## 8. Configuration

File: `config.yaml`

```yaml
check_interval_seconds: 30
request_timeout_seconds: 5

telegram:
  enabled: true
  bot_token: "REPLACE_ME"
  chat_id: "REPLACE_ME"

targets:
  - name: "Example"
    url: "https://example.com"
  - name: "Google"
    url: "https://google.com"
```

A `config.example.yaml` is committed. The real `config.yaml` is gitignored.

## 9. Concurrency Model

* The scheduler runs one check pass at a time.
* Inside a pass, `checker.CheckAll` spawns one goroutine per target.
* Results are collected through a buffered channel sized to `len(targets)`,
  so no goroutine blocks waiting for a reader.
* Each goroutine uses the same `*http.Client`, which is safe for concurrent use.
* Every request carries the shutdown `context`, so a SIGTERM cancels
  in-flight requests instead of leaving them hanging until timeout.
* `http.Client.Timeout` is a second safety net for requests that ignore
  context cancellation.

## 10. Error Handling

* Errors during startup (config, DB, Telegram) → log and `os.Exit(1)`.
* Errors during a check (timeout, DNS, TLS) → captured in `Result.Error`,
  status set to `DOWN`, never crashes the process.
* Errors during persistence or notification → logged via `slog.Error`,
  the loop continues.
* No `panic` in production paths. Panics are reserved for programmer errors
  (invariants that should never break).

## 11. Observability

* Structured logging with `log/slog` (text handler to stdout in V1).

* Every check emits one log line:

  `msg=check target=Example status=UP latency_ms=123 status_code=200`

* Every state transition emits a log line and a Telegram message.

* No metrics endpoint in V1 (planned for V2, see `roadmap.md`).

## 12. Anti-Spam Alert Logic

The notifier keeps an in-memory map keyed by URL storing the last known
status. A message is sent only when the status differs from the last known
one. This prevents a target that is continuously DOWN from generating one
message per check cycle.

The map is process-local. On restart, the first check for each target will
always emit a notification, which is acceptable for V1.

## 13. Testing Strategy

* **Unit tests** for `checker` using `httptest.NewServer` to simulate
  UP, DOWN, timeout, and slow responses.
* **Unit tests** for `storage` using an in-memory SQLite (`:memory:`)
  so no file is created during test runs.
* **Unit tests** for `notifier` with a fake `Notifier` implementation
  that records messages instead of sending them.
* **Scheduler tests** use a short interval (10ms) and assert results
  are forwarded.
* No integration tests in V1. No mocks framework — hand-written fakes
  only.
* Target coverage: 70% on `internal/`.

## 14. Makefile Targets

| Target       | Action                                  |
| ------------ | --------------------------------------- |
| `make run`   | `go run ./cmd/monitor`                  |
| `make build` | `go build -o bin/monitor ./cmd/monitor` |
| `make test`  | `go test ./...`                         |
| `make lint`  | `golangci-lint run`                     |
| `make tidy`  | `go mod tidy`                           |

## 15. Technical Decisions

| Decision            | Choice                          | Reason                                                                           | Trade-off                                               |
| ------------------- | ------------------------------- | -------------------------------------------------------------------------------- | ------------------------------------------------------- |
| Language            | Go 1.22+                        | Native concurrency, single static binary, predictable deployment and performance | Smaller ecosystem, no REPL, more verbose error handling |
| HTTP layer          | `net/http` stdlib               | Learn the fundamentals, no framework lock-in                                     | More manual code than with a router                     |
| Database            | SQLite via `modernc.org/sqlite` | Zero CGO, single-file DB, trivial deploy                                         | Not suitable for multi-node deployments                 |
| DB access           | `database/sql`                  | Explicit SQL, no ORM magic                                                       | More boilerplate                                        |
| Config format       | YAML                            | Human-readable, widely understood                                                | Requires runtime validation                             |
| Logging             | `log/slog`                      | Stdlib, structured, no external dependency                                       | Fewer features than Zap or Zerolog                      |
| Alert channel       | Telegram Bot API                | Free, simple HTTP API, low latency                                               | External dependency, per-chat rate limits               |
| Scheduling          | `time.Ticker`                   | Stdlib, sufficient for a fixed interval                                          | No cron expressions in V1                               |
| Signal handling     | `signal.NotifyContext`          | Stdlib, integrates cleanly with `context`                                        | None                                                    |
| Process supervision | systemd                         | Native on Linux VPS, simple unit file                                            | Linux-only                                              |

## 16. Directory Layout

```text
uptime-monitor/
├── cmd/
│   └── monitor/
│       └── main.go
├── internal/
│   ├── domain/
│   │   └── result.go
│   ├── config/
│   │   └── config.go
│   ├── checker/
│   │   ├── checker.go
│   │   └── checker_test.go
│   ├── storage/
│   │   ├── sqlite.go
│   │   ├── migrate.go
│   │   └── sqlite_test.go
│   ├── scheduler/
│   │   ├── scheduler.go
│   │   └── scheduler_test.go
│   └── notifier/
│       ├── notifier.go
│       └── telegram.go
├── migrations/
│   └── 001_init.sql
├── config.yaml
├── config.example.yaml
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
├── scope.md
├── architecture.md
├── roadmap.md
└── README.md
```

## 17. Extension Points

* **Swap SQLite for Postgres.** Add `internal/storage/postgres.go`, keep
  the `Store` interface unchanged, choose the implementation in `main`.

* **Add a new alert channel.** Add `internal/notifier/slack.go` implementing
  `Notifier`. Fan out in `main` by wrapping multiple notifiers.

* **Add content checks.** Extend `checker.Check` to assert on response body
  or headers. Add fields to `domain.Result` as needed.

* **Add a web dashboard.** Create `cmd/api` that reads the same SQLite file
  or points to Postgres. The monitor process stays unchanged.

* **Distributed monitoring.** Replace the in-process scheduler with a queue
  such as Redis and run multiple workers.

## 18. Deployment

Target: single VPS (Debian 12 or Ubuntu 24.04).

Steps:

1. Build for Linux:

   ```bash
   GOOS=linux GOARCH=amd64 go build -o bin/monitor ./cmd/monitor
   ```

2. Copy the binary, `config.yaml`, and `migrations/` to:

   ```text
   /opt/uptime-monitor/
   ```

3. Create a systemd unit at:

   ```text
   /etc/systemd/system/uptime-monitor.service
   ```

   ```ini
   [Unit]
   Description=Uptime Monitor
   After=network-online.target

   [Service]
   Type=simple
   WorkingDirectory=/opt/uptime-monitor
   ExecStart=/opt/uptime-monitor/bin/monitor
   Restart=on-failure
   RestartSec=5
   User=uptime

   [Install]
   WantedBy=multi-user.target
   ```

4. Enable and start the service:

   ```bash
   systemctl daemon-reload && systemctl enable --now uptime-monitor
   ```

5. Verify:

   ```bash
   journalctl -u uptime-monitor -f
   ```

## 19. Non-Goals for V1

* Multi-tenant authentication and billing.
* Distributed monitoring across regions.
* Sub-second resolution.
* Alert escalation policies.
* Public status page.

These are explicitly out of scope in V1. See `scope.md`.

## 20. Runtime Requirements

* Linux (or any OS with a Go toolchain for local development).
* Go 1.22+ for building.
* ~1 vCPU and ~50 MB RAM for the running binary.
* A writable directory for `uptime.db` and `config.yaml`.
* Outbound HTTPS access to monitored targets and to `api.telegram.org`.
* Estimated cost on a small VPS: ~$5/month (Hetzner CX22, Fly.io free tier,
  or similar).

---
