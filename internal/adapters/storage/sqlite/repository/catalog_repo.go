package repository

import (
	"autohost-cli/internal/domain"
	"autohost-cli/internal/ports"
	"context"
	"database/sql"
)

type catalogRepo struct{ db *sql.DB }

func NewCatalogRepo(db *sql.DB) ports.CatalogRepository { return &catalogRepo{db} }

func (r *catalogRepo) ListApps(ctx context.Context) ([]domain.CatalogApp, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT name, description
		FROM catalog_apps ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.CatalogApp
	for rows.Next() {
		var item domain.CatalogApp
		if err := rows.Scan(&item.Name, &item.Description); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *catalogRepo) FindByName(ctx context.Context, name domain.AppName) (domain.CatalogApp, error)
