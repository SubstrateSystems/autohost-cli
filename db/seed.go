// internal/db/seed.go
package db

import (
	"embed"
	"fmt"
	"sort"
)

//go:embed seeds/*.sql
var seedFS embed.FS

// Seed ejecuta los archivos de seeding en orden
func Seed(d *DB) error {
	// Tabla para trackear seeds ejecutados
	if _, err := d.Exec(`CREATE TABLE IF NOT EXISTS _seeds(id TEXT PRIMARY KEY);`); err != nil {
		return err
	}

	entries, err := seedFS.ReadDir("seeds")
	if err != nil {
		return err
	}

	// Ordena por nombre (001_, 002_...)
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })

	for _, e := range entries {
		var done string
		_ = d.QueryRow(`SELECT id FROM _seeds WHERE id=?`, e.Name()).Scan(&done)
		if done != "" {
			continue // Ya ejecutado
		}

		sqlBytes, err := seedFS.ReadFile("seeds/" + e.Name())
		if err != nil {
			return fmt.Errorf("error reading seed file %s: %w", e.Name(), err)
		}

		tx, err := d.Begin()
		if err != nil {
			return fmt.Errorf("error starting transaction for seed %s: %w", e.Name(), err)
		}

		if _, err := tx.Exec(string(sqlBytes)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error executing seed %s: %w", e.Name(), err)
		}

		if _, err := tx.Exec(`INSERT INTO _seeds(id) VALUES (?)`, e.Name()); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error marking seed %s as completed: %w", e.Name(), err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("error committing seed %s: %w", e.Name(), err)
		}
	}
	return nil
}
