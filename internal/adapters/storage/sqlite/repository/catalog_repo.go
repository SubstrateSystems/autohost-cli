package repository

import (
	"autohost-cli/internal/domain"
	"context"
	"database/sql"
)

type catalogRepo struct{ db *sql.DB }

func NewCatalogRepo(db *sql.DB) domain.CatalogRepo { return &catalogRepo{db} }

func (r *catalogRepo) ListApps(ctx context.Context) ([]domain.CatalogItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT name, description
		FROM catalog_apps ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.CatalogItem
	for rows.Next() {
		var item domain.CatalogItem
		if err := rows.Scan(&item.Name, &item.Description); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}
