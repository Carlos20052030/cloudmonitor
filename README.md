# CloudMonitor

A small Go service that checks HTTP endpoints on a schedule and
reports their status.

Built as a learning project with production habits: clean package
layout, dependency injection, and tests.

## Status

V1 — in development. See [docs/roadmap.md](./docs/roadmap.md) for the
current phase.

**Implemented:**

- Parallel HTTP checks with goroutines
- Configurable timeout via YAML
- UP/DOWN classification and latency measurement
- Unit tests for the checker using `httptest`

**Next:**

- SQLite persistence
- Telegram alerts on state transitions
- Scheduler with fixed interval
- Graceful shutdown

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

You should see one line per target, with status, HTTP code, and latency.

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

`check_interval_seconds` and the `telegram` block are placeholders for
features not yet implemented in V1. They are parsed and ignored today.

## Project Layout

```
cmd/monitor/        entry point
internal/domain/    shared types (Target, Result, Status)
internal/config/    YAML config loader
internal/checker/   HTTP checks with goroutines
docs/               scope, architecture, roadmap
```

## Documentation

- [docs/scope.md](./docs/scope.md) — what the project does and does not do (V1)
- [docs/architecture.md](./docs/architecture.md) — target architecture for V1
- [docs/roadmap.md](./docs/roadmap.md) — build plan and milestones

## License

MIT
```


