// internal/db/migrate.go
package db

import (
	"embed"
	"sort"
)

//go:embed migrations/*.sql
var migFS embed.FS

func Migrate(d *DB) error {
	if _, err := d.Exec(`CREATE TABLE IF NOT EXISTS _migrations(id TEXT PRIMARY KEY);`); err != nil {
		return err
	}

	entries, err := migFS.ReadDir("migrations")
	if err != nil {
		return err
	}
	// ordena por nombre (001_, 002_...)
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })

	for _, e := range entries {
		var done string
		_ = d.QueryRow(`SELECT id FROM _migrations WHERE id=?`, e.Name()).Scan(&done)
		if done != "" {
			continue
		}

		sqlBytes, err := migFS.ReadFile("migrations/" + e.Name())
		if err != nil {
			return err
		}

		tx, err := d.Begin()
		if err != nil {
			return err
		}
		if _, err := tx.Exec(string(sqlBytes)); err != nil {
			_ = tx.Rollback()
			return err
		}
		if _, err := tx.Exec(`INSERT INTO _migrations(id) VALUES (?)`, e.Name()); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}
