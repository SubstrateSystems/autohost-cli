package repository

import (
	"autohost-cli/internal/adapters/storage/sqlite/mappers"
	"autohost-cli/internal/adapters/storage/sqlite/models"
	"autohost-cli/internal/domain"
	"autohost-cli/internal/ports"
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type catalogRepo struct{ db *sqlx.DB }

func NewCatalogRepo(db *sqlx.DB) ports.CatalogRepository { return &catalogRepo{db} }

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

func (r *catalogRepo) FindByName(ctx context.Context, name domain.AppName) (domain.CatalogApp, error) {
	var row models.CatalogAppRow
	err := r.db.GetContext(ctx, &row, `
        SELECT id, name, description, default_port, default_port_db, client_db, created_at, updated_at
        FROM catalog_apps
        WHERE name = ?`, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.CatalogApp{}, nil
		}
		return domain.CatalogApp{}, err
	}
	return mappers.ToDomainCatalogApp(row), nil
}
