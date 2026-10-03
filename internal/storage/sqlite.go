package storage

import (
	"context"
	"database/sql"

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
