# Uptime Monitor

A small Go service that checks HTTP endpoints on a schedule, stores every
result in SQLite, and alerts a Telegram chat when a target changes state
(UP → DOWN or DOWN → UP).

Built as a learning project with production habits: clean package layout,
dependency injection, graceful shutdown, structured logs, and tests.

## Features

- Parallel HTTP checks using goroutines
- Configurable check interval and request timeout
- SQLite persistence with indexed history
- Telegram alerts on state transitions (no spam)
- Graceful shutdown on SIGINT / SIGTERM
- Structured logging with `log/slog`
- Zero CGO, single static binary

## Stack

- Go 1.22+
- `net/http` (standard library)
- `database/sql` + `modernc.org/sqlite`
- `log/slog` (standard library)
- `gopkg.in/yaml.v3`
- Telegram Bot API

No web framework. No ORM. Standard library first.

## Quick Start

1. Clone the repository:

   ```bash
   git clone https://github.com/Carlos20052030/cloudmonitor.git
   cd cloudmonitor
   ```

2. Copy the example configuration and edit it:

   ```bash
   cp config.example.yaml config.yaml
   ```

   Fill in your Telegram `bot_token` and `chat_id`, and list the URLs you
   want to monitor.

3. Run it:

   ```bash
   make run
   ```

   You should see one log line per check in the terminal, and a Telegram
   message when a target changes state.

## Configuration

`config.yaml` controls everything:

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

See `config.example.yaml` for a template. The real `config.yaml` is
gitignored.

## Make Targets

| Target       | Action                                  |
| ------------ | --------------------------------------- |
| `make run`   | Run the monitor locally                 |
| `make build` | Build the binary into `bin/monitor`     |
| `make test`  | Run all tests                           |
| `make lint`  | Run `golangci-lint`                     |
| `make tidy`  | Tidy `go.mod` / `go.sum`                |

## Project Layout

```
cmd/monitor/        entry point and dependency wiring
internal/domain/    shared types (Target, Result, Status)
internal/config/    YAML config loader
internal/checker/   HTTP checks with goroutines
internal/storage/   SQLite persistence
internal/notifier/  Telegram alerts
internal/scheduler/ periodic check runner
migrations/         SQL migrations
```

## Documentation

- [scope.md](./scope.md) — what the project does and does not do (V1)
- [architecture.md](./architecture.md) — how it is built
- [roadmap.md](./roadmap.md) — build plan, milestones, and risks

## Status

V1 — in development. See [roadmap.md](./roadmap.md) for the current phase.

## License

MIT
```

---