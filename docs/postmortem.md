# Postmortem — CloudMonitor V1

A retrospective on building CloudMonitor V1: what worked, what didn't,
and what I would do differently.

## Context

CloudMonitor started as a learning project with a market goal: build a
real uptime monitor in Go, self-hosted, with Telegram alerts, and ship
it as part of a backend portfolio. The plan was 4 weeks, one commit per
day, four phases.

## What worked

**Small commits, every day.** The habit of committing daily — even when
the change was one line — kept momentum. The git log became a progress
log. Every commit message describes what changed and why.

**Test-first for hard parts.** The checker was tested with `httptest`
before being wired into the main. That caught the timeout behavior early
and made `CheckAll` safe to build on top of.

**Interfaces at the boundary.** `Store` and `Notifier` are declared
where they are consumed, not where they are implemented. Swapping SQLite
for Postgres later would be a new file, not a rewrite.

**`//go:embed` for the SQL schema.** The migration file is embedded in
the binary. Tests stopped breaking when the working directory changed.

**Anti-spam from day one.** The notifier only sends messages on state
transitions. Without this, the bot becomes unusable in minutes.

## What didn't work

**`os.ReadFile` for migrations.** Worked locally, broke in tests because
Go changes the working directory when running tests per package.
Replaced with `//go:embed`.

**Config with relative paths.** `config.yaml` is read relative to the
process working directory. Fine for now, but a systemd deployment needs
an absolute path or a `WorkingDirectory=` directive.

**No integration test for the scheduler with a real HTTP server.** The
current tests use `httptest`, which is good, but the full loop
(check → save → notify) is only tested manually.

## What I learned

- Go's `context` threading is worth doing from day one. Adding it later
  is painful.
- `signal.NotifyContext` is the cleanest way to handle shutdown. The
  whole graceful shutdown is three lines.
- Structured logs (`log/slog`) beat `fmt.Println` in every dimension
  once the program has more than one goroutine.
- Telegram bot tokens leak easily. Revoking in BotFather is the only
  safe response. Never commit `config.yaml`.
- `golangci-lint` catches real issues that the compiler ignores. The
  `errcheck` rule alone found six places where errors were silently
  dropped.

## What I would do differently

- Add a `Makefile` on day one, not week four. Running `make test` is
  faster to type than `go test ./...`.
- Write the `README.md` skeleton early and fill it in as features land,
  instead of writing it at the end.
- Test the notifier's Telegram integration with a fake HTTP server
  instead of relying on manual verification.

## Metrics

- Planned as a 4-week project; executed in focused sessions over 6 days.
- 15 unit tests across `checker`, `storage`, `notifier`, and `scheduler`.
- 0 issues reported by `golangci-lint`.
- 1 binary, ~17 MB, statically linked, no CGO.
- 1 SQLite database file, ~16 KB after first run.

## Next

V1 is feature-complete for local use. V2 candidates: Prometheus metrics,
a web dashboard, additional alert channels, and Postgres support. See
`roadmap.md` for the full list.