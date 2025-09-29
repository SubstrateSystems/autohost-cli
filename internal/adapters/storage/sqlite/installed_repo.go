package sqlite

import (
	"context"
	"database/sql"

	"autohost-cli/internal/domain"
)

type InstalledRepo struct {
	db *sql.DB
}

func NewInstalledRepo(db *sql.DB) *InstalledRepo {
	return &InstalledRepo{db: db}
}

func (r *InstalledRepo) List(ctx context.Context) ([]domain.InstalledApp, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, created_at
		FROM installed_apps
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.InstalledApp
	for rows.Next() {
		var a domain.InstalledApp
		if err := rows.Scan(&a.ID, &a.Name, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *InstalledRepo) Add(ctx context.Context, app domain.InstalledApp) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO installed_apps (name, catalog_app_id) 
		VALUES (?, ?)
	`, app.Name, app.CatalogAppID)
	return err
}

func (r *InstalledRepo) Remove(ctx context.Context, name string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM installed_apps WHERE name = ?
	`, name)
	return err
}

func (r *InstalledRepo) IsInstalledApp(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM installed_apps WHERE name = ?
		)`, name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
