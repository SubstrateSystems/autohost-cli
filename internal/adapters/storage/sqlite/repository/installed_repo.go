package repository

import (
	"context"
	"database/sql"
	"time"

	"autohost-cli/internal/domain"
	"autohost-cli/internal/ports"

	"github.com/jmoiron/sqlx"
)

type installedRepo struct {
	db *sqlx.DB
}

func NewInstalledRepo(db *sqlx.DB) ports.InstalledRepository {
	return &installedRepo{db: db}
}

func (r *installedRepo) List(ctx context.Context) ([]domain.InstalledApp, error) {
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
		var (
			id   int64
			name string
			ts   sql.NullInt64
		)
		if err := rows.Scan(&id, &name, &ts); err != nil {
			return nil, err
		}

		var createdAt time.Time
		if ts.Valid {
			createdAt = time.Unix(ts.Int64, 0).UTC()
		}

		out = append(out, domain.InstalledApp{
			ID:        int64(id),
			Name:      name,
			CreatedAt: createdAt,
		})
	}
	return out, rows.Err()
}

func (r *installedRepo) Install(ctx context.Context, app domain.InstalledApp) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO installed_apps (name, port, port_db, catalog_app_id) 
		VALUES (?, ?, ?, ?)
	`, app.Name, app.Port, app.PortDB, app.CatalogAppID)
	return err
}

func (r *installedRepo) Remove(ctx context.Context, name domain.AppName) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM installed_apps WHERE name = ?
	`, name)
	return err
}

func (r *installedRepo) IsInstalled(ctx context.Context, name domain.AppName) (bool, error) {
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
