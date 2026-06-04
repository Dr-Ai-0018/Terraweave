package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct {
	db *sql.DB
}

type Status struct {
	OK            bool   `json:"ok"`
	Postgres      string `json:"postgres,omitempty"`
	PostGIS       string `json:"postgis,omitempty"`
	CurrentUser   string `json:"current_user,omitempty"`
	CurrentSchema string `json:"current_schema,omitempty"`
	Error         string `json:"error,omitempty"`
}

func Open(databaseURL string) (*Store, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("missing TERRAWEAVE_DATABASE_URL")
	}
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) Check(ctx context.Context) Status {
	if s == nil || s.db == nil {
		return Status{OK: false, Error: "database not configured"}
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var status Status
	err := s.db.QueryRowContext(ctx, `
		select
			version(),
			postgis_version(),
			current_user,
			current_schema()
	`).Scan(&status.Postgres, &status.PostGIS, &status.CurrentUser, &status.CurrentSchema)
	if err != nil {
		status.Error = err.Error()
		return status
	}
	status.OK = true
	return status
}
