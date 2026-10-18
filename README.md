# CloudMonitor

![Go Version](https://img.shields.io/badge/Go-1.27-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)
![Status](https://img.shields.io/badge/status-live-brightgreen)

**A production-ready uptime monitor for HTTP endpoints, built with Go.**

## The Problem

Downtime is expensive. You need to know the second your service goes offline, but you also don't want to be spammed with alerts every 30 seconds just because a network blip happened.

CloudMonitor solves this. It checks your endpoints on a schedule in parallel, measures latency, stores the history in SQLite, and sends a Telegram alert **only when the state changes** (from UP to DOWN, or DOWN to UP).

It runs 24/7 on a Linux VPS with zero external dependencies.

## Features

- **Parallel HTTP Checks:** Built with goroutines for speed and efficiency.
- **Smart Alerting:** Telegram notifications are sent *only* on state transitions (anti-spam).
- **Persistent History:** SQLite storage with indexed queries for uptime history and latency tracking.
- **Zero Framework:** Pure standard library `net/http`, `database/sql`, and `log/slog`.
- **Graceful Shutdown:** In-flight checks finish before the process exits (handles `SIGINT`/`SIGTERM`).
- **Deployable:** Designed to run under `systemd` on a Linux VPS.

## Tech Stack

| Layer | Choice |
|---|---|
| Language | Go 1.27 |
| HTTP | `net/http` (stdlib) |
| Database | SQLite (`modernc.org/sqlite`) |
| Logging | `log/slog` (structured) |
| Migrations | Embedded `.sql` files (`//go:embed`) |

## Quick Start

### Prerequisites
- Go 1.27+

### 1. Clone and Configure
```bash
git clone https://github.com/Carlos20052030/cloudmonitor.git
cd cloudmonitor
cp config.example.yaml config.yaml
```

### 2. Edit the Config
List the URLs you want to monitor in `config.yaml`:
```yaml
check_interval_seconds: 30
request_timeout_seconds: 5

telegram:
  enabled: false # Change to true if you set up a bot

targets:
  - name: "Example"
    url: "https://example.com"
  - name: "Google"
    url: "https://google.com"
```

### 3. Run it
```bash
go run ./cmd/monitor
```
You will see one log line per target, with status, HTTP code, and latency. Results are saved to `cloudmonitor.db`.

## Telegram Setup

To receive alerts, you need a Telegram bot:

1. Open Telegram and talk to [@BotFather](https://t.me/BotFather).
2. Send `/newbot` and follow the prompts.
3. Copy the token BotFather gives you.
4. Open a chat with your new bot and send `/start`.
5. Access `https://api.telegram.org/bot<YOUR_TOKEN>/getUpdates` in a browser to find your `chat.id`.
6. Update `config.yaml`:
   ```yaml
   telegram:
     enabled: true
     bot_token: "YOUR_TOKEN"
     chat_id: "123456789"
   ```

## Project Layout

```text
cmd/monitor/        entry point and dependency wiring
internal/domain/    shared types (Target, Result, Status)
internal/config/    YAML config loader
internal/checker/   HTTP checks with goroutines
internal/storage/   SQLite persistence
internal/scheduler/ periodic check runner
internal/notifier/  Telegram alerts with anti-spam
migrations/         SQL schema, embedded via //go:embed
docs/               scope, architecture, roadmap, postmortem
```

## Documentation

- [docs/scope.md](./docs/scope.md) — problem, decisions, success criteria
- [docs/architecture.md](./docs/architecture.md) — architecture and concurrency model
- [docs/deployment.md](./docs/deployment.md) — how to deploy on a VPS with systemd
- [docs/postmortem.md](./docs/postmortem.md) — retrospective on building V1

## License

MIT — see [LICENSE](LICENSE).
```