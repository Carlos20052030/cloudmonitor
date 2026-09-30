# Architecture

## Overview
CloudMonitor v0.1 is a single-process Python application.
It loads a list of HTTP endpoints, checks them periodically,
stores results in SQLite, and serves a simple status page.

## Components
- **Config Loader**: reads `config.yaml`.
- **Scheduler**: triggers checks at a fixed interval.
- **Checker**: performs HTTP requests and measures latency.
- **Storage**: persists results in SQLite.
- **Web**: exposes a status page and JSON API.

## Data Flow
1. The application starts.
2. Config is loaded.
3. The scheduler triggers checks.
4. The checker performs HTTP requests.
5. Results are saved to SQLite.
6. The web server reads the latest results.

## Storage Schema
Table `checks`:
- `id`
- `url`
- `status` (up/down)
- `status_code`
- `latency_ms`
- `checked_at`

## Tech Stack
- Python 3.11+
- FastAPI or Flask
- SQLite
- httpx or requests
- PyYAML

## Constraints
- Single process
- Local only
- No external database
- No alerting yet