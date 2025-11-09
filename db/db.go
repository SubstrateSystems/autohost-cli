// internal/db/db.go
package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type DB struct{ *sqlx.DB }

func Open(path string) (*DB, error) {
	// WAL + busy_timeout + FK
	dsn := fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", path)
	raw, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if _, err := raw.Exec(`PRAGMA foreign_keys=ON;`); err != nil {
		return nil, err
	}
	return &DB{raw}, nil
}
