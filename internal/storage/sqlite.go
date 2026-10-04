package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/Carlos20052030/cloudmonitor/internal/domain"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func New(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) Save(ctx context.Context, r domain.Result) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO checks
			(target_name, url, status, status_code, latency_ms, error, checked_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, r.TargetName, r.URL, r.Status, r.StatusCode, r.LatencyMS, r.Error, r.CheckedAt)
	return err
}

func (s *Store) Close() error {
	return s.db.Close()
}

type CheckRecord struct {
	ID         int64
	TargetName string
	URL        string
	Status     domain.Status
	StatusCode int
	LatencyMS  int64
	Error      string
	CheckedAt  time.Time
}

// RecentChecks returns the last n checks for a target, most recent first.
func (s *Store) RecentChecks(ctx context.Context, targetName string, n int) ([]CheckRecord, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, target_name, url, status, status_code, latency_ms, error, checked_at
		FROM checks
		WHERE target_name = ?
		ORDER BY checked_at DESC
		LIMIT ?
	`, targetName, n)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var out []CheckRecord
	for rows.Next() {
		var r CheckRecord
		if err := rows.Scan(
			&r.ID, &r.TargetName, &r.URL, &r.Status,
			&r.StatusCode, &r.LatencyMS, &r.Error, &r.CheckedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, r)
	}

	return out, rows.Err()
}
