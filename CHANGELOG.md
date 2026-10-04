# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

Nothing yet.

## [1.0.0] - 2026-10-04

First stable release of CloudMonitor V1.

### Added

- Parallel HTTP checks using goroutines and a buffered channel.
- YAML configuration for interval, timeout, Telegram, and targets.
- UP/DOWN classification based on HTTP status code (2xx = UP).
- Latency measurement per check, in milliseconds.
- SQLite persistence for every check result, with an index on
  `(target_name, checked_at)`.
- `RecentChecks` query to retrieve the last N checks for a given target.
- Telegram notifier with anti-spam: one message per state transition.
- Scheduler with `time.Ticker` running one check pass at a fixed interval.
- Graceful shutdown on SIGINT and SIGTERM via `signal.NotifyContext`.
- Structured logging with `log/slog`.
- `//go:embed` for the SQL schema, so the binary is self-contained.
- Unit tests for `checker`, `storage`, `notifier`, and `scheduler`.
- `Makefile` with `run`, `build`, `test`, `lint`, `tidy`, and `clean`.
- `README.md` with badges, demo GIF, quick start, and Telegram setup.
- `LICENSE` (MIT).
- Project documentation: `scope.md`, `architecture.md`, `roadmap.md`,
  `postmortem.md`.

### Notes

- The VPS deployment (systemd unit) is documented in `architecture.md`
  but not yet executed. It is the only remaining item of the original
  roadmap.