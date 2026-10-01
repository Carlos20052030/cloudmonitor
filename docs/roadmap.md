# Roadmap — Uptime Monitor V1

Build plan for V1. Four weeks, one commit per day, four phases.

## Overview

- **Duration:** 4 weeks (28 days)
- **Cadence:** at least one commit per day, even if small
- **Definition of done:** `go test ./...` passes, `golangci-lint run` is clean,
  the binary runs on a VPS for 7 consecutive days without crashing, and the
  README lets a stranger run the project in 10 minutes.
- **Working rule:** do not start the next phase until the previous one runs.

## Phase 1 — Foundation (Week 1)

**Goal:** the project runs and checks multiple URLs in parallel.

**Deliverable:** a binary that checks every target in the config and prints
one log line per target.

**Done when:** `go test ./...` passes and three URLs are checked concurrently.

| Day | Task                                                                 |
| --- | -------------------------------------------------------------------- |
| 1   | `go mod init`, directory layout, `.gitignore`, empty `README.md`     |
| 2   | `internal/domain/result.go` with `Target`, `Result`, `Status`        |
| 3   | `internal/config/config.go` + `config.yaml` + `config.example.yaml`  |
| 4   | `internal/checker/checker.go` — check a single URL                   |
| 5   | `checker.CheckAll` using goroutines and a buffered channel           |
| 6   | `internal/checker/checker_test.go` using `httptest.NewServer`        |
| 7   | Commit, update `README.md`, push to GitHub                           |

## Phase 2 — Persistence (Week 2)

**Goal:** every check is stored in SQLite.

**Deliverable:** the monitor writes to `uptime.db` and history is queryable.

**Done when:** `sqlite3 uptime.db "SELECT * FROM checks"` returns rows from
the last run.

| Day | Task                                                                 |
| --- | -------------------------------------------------------------------- |
| 8   | `migrations/001_init.sql` + `internal/storage/migrate.go`            |
| 9   | `internal/storage/sqlite.go` with `Save` and `Close`                 |
| 10  | Wire storage into `cmd/monitor/main.go`                              |
| 11  | `internal/storage/sqlite_test.go` using `:memory:` SQLite            |
| 12  | Index on `(target_name, checked_at)` and a "last N per target" query |
| 13  | Define `Store` interface in `main`, inject it                        |
| 14  | Commit, document the schema in `architecture.md` if changed          |

## Phase 3 — Automation and Alerts (Week 3)

**Goal:** the service runs on its own and alerts on state changes.

**Deliverable:** a long-running process that checks on an interval and sends
Telegram messages only when a target changes state.

**Done when:** manually taking a target down produces exactly one Telegram
message, and restoring it produces another.

| Day | Task                                                                 |
| --- | -------------------------------------------------------------------- |
| 15  | `internal/scheduler/scheduler.go` with `time.Ticker`                 |
| 16  | Graceful shutdown via `signal.NotifyContext` and `context`           |
| 17  | `internal/notifier/notifier.go` interface + `telegram.go`            |
| 18  | Anti-spam logic: notify only on state transition                     |
| 19  | `scheduler_test.go` and `notifier_test.go` with hand-written fakes   |
| 20  | Structured logging with `log/slog` on every check and transition     |
| 21  | Commit, update `README.md` with Telegram setup steps                 |

## Phase 4 — Production (Week 4)

**Goal:** running 24/7 on a VPS, portfolio-ready.

**Deliverable:** systemd service active on a VPS, README finished, GIF of
the Telegram alert attached.

**Done when:** the service has been running for 7 days without manual
intervention and a stranger can reproduce the setup from the README.

| Day | Task                                                                 |
| --- | -------------------------------------------------------------------- |
| 22  | `Makefile` with `run`, `build`, `test`, `lint`, `tidy`               |
| 23  | Deploy on a VPS with the systemd unit from `architecture.md`         |
| 24  | End-to-end test: kill a real target, confirm the alert arrives       |
| 25  | Polish `README.md`: badges, quick start, screenshot                  |
| 26  | `golangci-lint run` clean, address every warning                     |
| 27  | Record a short GIF showing the Telegram alert, embed in `README.md`  |
| 28  | Publish on LinkedIn, pin the repo, write a short post-mortem         |

## Milestones

| Milestone | End of      | What it proves                                      |
| --------- | ----------- | --------------------------------------------------- |
| M1        | Week 1      | Parallel checks work, tests pass                    |
| M2        | Week 2      | Data is persisted and queryable                     |
| M3        | Week 3      | The service runs unattended and notifies correctly  |
| M4        | Week 4      | The service survives in production for 7 days       |

## Risks

| Risk                                          | Mitigation                                                    |
| --------------------------------------------- | ------------------------------------------------------------- |
| Stuck on goroutines or channels               | Re-read Go concurrency basics, use a buffered channel, ask for help after 2 hours |
| Cannot deploy to a VPS                        | Fall back to Fly.io or a free tier host                       |
| Losing time polishing a dashboard that is out of scope | Dashboard is a non-goal for V1. Stop and return to the plan. |
| Losing motivation mid-project                 | One commit per day, even if it is a comment fix               |
| Telegram bot token leaking in the repo        | Keep `config.yaml` in `.gitignore` from day 1                 |
| Overengineering storage before it is needed   | Ship SQLite first; Postgres is a V2 concern                   |

## Post-V1 (not now)

- Web dashboard (`cmd/api` reading the same database)
- Additional alert channels: Slack, Discord, SMS, WhatsApp
- Content checks: assert on response body or headers
- Multi-user support with authentication and plans
- Distributed monitoring across regions
- Public status page
- Metrics endpoint (Prometheus) and Grafana dashboards
- Cron expressions instead of a fixed interval