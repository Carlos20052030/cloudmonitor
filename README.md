# CloudMonitor

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-blue)](./LICENSE)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen)](#)
[![Lint](https://img.shields.io/badge/golangci--lint-clean-brightgreen)](#)
[![Status](https://img.shields.io/badge/status-V1-yellow)](./docs/roadmap.md)

A small Go service that checks HTTP endpoints on a schedule and
reports their status.

Built as a learning project with production habits: clean package
layout, dependency injection, and tests.

## Status

V1 — in development. See [docs/roadmap.md](./docs/roadmap.md) for the
current phase.

**Implemented:**

- Parallel HTTP checks with goroutines
- Configurable check interval and request timeout via YAML
- UP/DOWN classification and latency measurement
- SQLite persistence with indexed history
- Telegram alerts only on state transitions (no spam)
- Scheduler with fixed interval
- Graceful shutdown on SIGINT / SIGTERM
- Structured logging with `log/slog`
- Unit tests for checker, storage, notifier, and scheduler

**Next:**

- `Makefile` with `run`, `build`, `test`, `lint`, `tidy`
- Deployment under systemd on a Linux VPS
- `golangci-lint` clean run

## Quick Start

1. Clone the repository:

   ```bash
   git clone https://github.com/Carlos20052030/cloudmonitor.git
   cd cloudmonitor
   ```

2. Copy the example configuration:

   ```bash
   cp config.example.yaml config.yaml
   ```

3. Edit `config.yaml` and list the URLs you want to monitor.

4. Run it:

   ```bash
   go run ./cmd/monitor
   ```

You should see one log line per target, with status, HTTP code, and latency.
Results are saved to `cloudmonitor.db`.

## Demo

```
time=2026-10-04T07:38:46-03:00 level=INFO msg="telegram enabled" chat_id=123456789
time=2026-10-04T07:38:46-03:00 level=INFO msg="monitor starting" targets=3 interval=30s
time=2026-10-04T07:38:46-03:00 level=INFO msg=check target="Example" status=UP status_code=200 latency_ms=238
time=2026-10-04T07:38:46-03:00 level=INFO msg=check target="Google" status=UP status_code=200 latency_ms=519
time=2026-10-04T07:38:46-03:00 level=INFO msg=check target="Site Que Não Existe" status=DOWN status_code=0 latency_ms=206
^C
time=2026-10-04T07:39:00-03:00 level=INFO msg="monitor stopped"
```

When a target changes state, a Telegram message is sent:

```
🔴 DOWN: Site Que Não Existe
URL: https://site-que-nao-existe-12345.com
Error: dial tcp: lookup site-que-nao-existe-12345.com: no such host
```

```
🟢 UP: Example
URL: https://example.com
Latency: 238ms
```

## Configuration

`config.yaml` controls everything:

```yaml
check_interval_seconds: 30
request_timeout_seconds: 5

telegram:
  enabled: false
  bot_token: "REPLACE_ME"
  chat_id: "REPLACE_ME"

targets:
  - name: "Example"
    url: "https://example.com"
  - name: "Google"
    url: "https://google.com"
```

See `config.example.yaml` for a template. The real `config.yaml` is
gitignored.

## Telegram Setup

Telegram alerts are optional. To enable them:

1. Open Telegram and talk to [@BotFather](https://t.me/BotFather).
2. Send `/newbot` and follow the prompts to create a bot.
3. Copy the token BotFather gives you.
4. Open a chat with your new bot and send `/start`. This is required —
   the bot cannot message you until you message it first.
5. Open `https://api.telegram.org/bot<YOUR_TOKEN>/getUpdates` in a browser.
6. Find the `chat.id` field in the JSON response. It looks like `123456789`.
7. Edit `config.yaml`:

   ```yaml
   telegram:
     enabled: true
     bot_token: "YOUR_TOKEN"
     chat_id: "123456789"
   ```

8. Run `go run ./cmd/monitor`. You should receive one message per target
   on the first check, and one message only when a target changes state.

**Security:** never commit `config.yaml`. It contains your bot token and
is already covered by `.gitignore`. If a token leaks, revoke it in
BotFather and generate a new one.

## Project Layout

```
cmd/monitor/        entry point and dependency wiring
internal/domain/    shared types (Target, Result, Status)
internal/config/    YAML config loader
internal/checker/   HTTP checks with goroutines
internal/storage/   SQLite persistence
internal/scheduler/ periodic check runner
internal/notifier/  Telegram alerts with anti-spam
migrations/         SQL schema, embedded via //go:embed
docs/               scope, architecture, roadmap
```

## Documentation

- [docs/scope.md](./docs/scope.md) — what the project does and does not do (V1)
- [docs/architecture.md](./docs/architecture.md) — target architecture for V1
- [docs/roadmap.md](./docs/roadmap.md) — build plan and milestones

## License

MIT
