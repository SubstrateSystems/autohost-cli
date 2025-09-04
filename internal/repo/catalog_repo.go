// internal/repo/catalog_repo.go
package repo

import (
	"autohost-cli/internal/models"
	"context"
	"database/sql"
)

type CatalogRepo interface {
	ListApps(ctx context.Context) ([]models.InstalledApp, error)
	// VersionsByApp(ctx context.Context, app string) ([]domain.CatalogVersion, error)
	// GetApp(ctx context.Context, name string) (*domain.CatalogApp, error)
	// UpsertApp(ctx context.Context, app domain.CatalogApp) error
}

type catalogRepo struct{ db *sql.DB }

func NewCatalogRepo(db *sql.DB) CatalogRepo { return &catalogRepo{db} }

func (r *catalogRepo) ListApps(ctx context.Context) ([]models.InstalledApp, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT name
		FROM catalog_apps ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.InstalledApp
	for rows.Next() {
		var model models.InstalledApp
		if err := rows.Scan(&model.Name); err != nil {
			return nil, err
		}
		out = append(out, model)
	}
	return out, rows.Err()
}
