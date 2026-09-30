# Scope

## Purpose
CloudMonitor v0.1 is a minimal HTTP/HTTPS uptime monitor.
It checks endpoints periodically and shows whether they are up or down.

## In Scope
- HTTP and HTTPS endpoint checks
- Configurable interval
- Status and latency recording
- Simple status page and JSON API
- Local execution
- Basic tests

## Out of Scope
- Alerting
- Authentication
- Multi-user support
- Cloud provider integrations
- System metrics (CPU, RAM, disk)
- Distributed collection
- Advanced dashboards

## Success Criteria
- A user can configure endpoints in YAML.
- The tool checks them automatically.
- Results are stored and visible.
- The project runs locally with one command.