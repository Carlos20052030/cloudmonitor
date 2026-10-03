# Scope — CloudMonitor V1

This document defines what V1 does and, more importantly, what V1 does not do.
Scope is the contract that prevents the project from growing forever.

## Problem

Developers and small teams usually discover that a website or API went down
when a customer complains. There is no cheap, self-hosted way to watch a
handful of HTTP endpoints and get a real alert when one of them flips.

## Objective

Watch a list of HTTP targets on a fixed interval, record every result, and
send a Telegram message when a target changes state (UP → DOWN or DOWN → UP).

## In Scope (V1)

- Check multiple HTTP targets in parallel using goroutines
- Configurable interval and request timeout via `config.yaml`
- HTTP status codes in the 2xx range count as UP; everything else counts as DOWN
- Persist every check to SQLite with an indexed history
- Send Telegram alerts only on state transitions (no spam)
- Graceful shutdown on SIGINT / SIGTERM
- Structured logging to stdout with `log/slog`
- Run as a single binary under systemd on a Linux VPS
- Tests for `checker`, `storage`, and `scheduler`
- A `Makefile` with `run`, `build`, `test`, `lint`, and `tidy`

## Out of Scope (V1)

These are explicitly deferred. Do not build them in V1.

- Web dashboard or any user interface
- User accounts, authentication, or multi-tenancy
- Billing, plans, or any commercial feature
- Additional alert channels (Slack, Discord, SMS, WhatsApp, e-mail, webhooks)
- Content checks (response body, headers, JSON shape)
- TLS certificate expiration checks
- DNS resolution checks
- Distributed monitoring from multiple regions
- Sub-second or per-second resolution
- Cron expressions; V1 uses a single fixed interval
- Alert escalation policies or on-call rotations
- Public status page
- Metrics endpoint (Prometheus) or Grafana dashboards
- Database other than SQLite
- Configuration hot reload; restart is required to apply changes
- Retry logic with backoff; a failed check is a failed check

## User

A solo developer or a small agency that runs a handful of sites or APIs and
wants a simple, self-hosted uptime monitor with Telegram alerts. The user is
comfortable editing a YAML file and deploying a Go binary to a VPS.

## Success Criteria

V1 is considered successful when all of the following are true:

1. `go test ./...` passes without failures.
2. `golangci-lint run` reports zero issues.
3. The service runs on a VPS for 7 consecutive days without crashing.
4. A target taken down manually produces exactly one Telegram alert within
   two check cycles.
5. Bringing the same target back up produces exactly one recovery alert.
6. A stranger can clone the repository and run the project in under 10
   minutes by following the `README.md`.
7. No secret (bot token, chat id) is ever committed to the repository.

## Constraints

- Build time: 4 weeks, one commit per day.
- No paid third-party services. Telegram and the VPS are the only external
  dependencies.
- Target hardware: 1 vCPU, 1 GB RAM VPS.
- Language: Go 1.22+.
- Third-party dependencies limited to `modernc.org/sqlite`,
  `gopkg.in/yaml.v3`, and `go-telegram-bot-api/v5`.
- No web framework, no ORM, no mocking framework.

## Assumptions

- The user has outbound HTTPS access to the monitored targets and to
  `api.telegram.org`.
- The user can create a Telegram bot and obtain a `chat_id`.
- The user can deploy a static Go binary and a systemd unit on a Linux VPS.

## Version

V1 — minimum viable product. No user interface. No commercial features.
Single process, single database, single alert channel.