package sqlite

import (
	"autohost-cli/internal/domain"
	"context"
	"database/sql"
)

type CatalogRepo interface {
	ListApps(ctx context.Context) ([]domain.CatalogApp, error)
	// VersionsByApp(ctx context.Context, app string) ([]domain.CatalogVersion, error)
	// GetApp(ctx context.Context, name string) (*domain.CatalogApp, error)
	// UpsertApp(ctx context.Context, app domain.CatalogApp) error
}

type catalogRepo struct{ db *sql.DB }

func NewCatalogRepo(db *sql.DB) CatalogRepo { return &catalogRepo{db} }

func (r *catalogRepo) ListApps(ctx context.Context) ([]domain.CatalogApp, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT name, description, created_at, updated_at
		FROM catalog_apps ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.CatalogApp
	for rows.Next() {
		var model domain.CatalogApp
		if err := rows.Scan(&model.Name, &model.Description, &model.CreatedAt, &model.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, model)
	}
	return out, rows.Err()
}
