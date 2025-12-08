package db

import (
	"database/sql"
	"fmt"

	"github.com/garrettladley/tsks/internal/migrations"
	"github.com/garrettladley/tsks/internal/sqlc"
	_ "github.com/mattn/go-sqlite3"
)

// Open opens a connection to the SQLite database and returns the queries instance.
// It automatically applies any pending migrations.
// The caller is responsible for closing the returned *sql.DB.
func Open(dbPath string) (*sql.DB, *sqlc.Queries, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := migrations.Apply(db, "migrations"); err != nil {
		db.Close()
		return nil, nil, err
	}

	queries := sqlc.New(db)
	return db, queries, nil
}
