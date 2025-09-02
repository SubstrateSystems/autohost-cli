// internal/repo/installed.go
package repo

import (
	"context"
	"database/sql"

	"autohost-cli/internal/models"
)

type InstalledRepo struct {
	db *sql.DB
}

func NewInstalledRepo(db *sql.DB) *InstalledRepo {
	return &InstalledRepo{db: db}
}

func (r *InstalledRepo) List(ctx context.Context) ([]models.InstalledApp, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, created_at
		FROM installed_apps
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.InstalledApp
	for rows.Next() {
		var a models.InstalledApp
		if err := rows.Scan(&a.ID, &a.Name, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *InstalledRepo) Add(ctx context.Context, app models.InstalledApp) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO installed_apps (name, catalog_app_id) 
		VALUES (?, ?)
	`, app.Name, app.CatalogAppID)
	return err
}
