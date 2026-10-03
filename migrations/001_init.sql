CREATE TABLE IF NOT EXISTS checks (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    target_name  TEXT     NOT NULL,
    url          TEXT     NOT NULL,
    status       TEXT     NOT NULL,
    status_code  INTEGER,
    latency_ms   INTEGER,
    error        TEXT,
    checked_at   DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_checks_target_time
    ON checks (target_name, checked_at);