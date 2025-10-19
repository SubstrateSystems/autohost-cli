package models

import "time"

type InstalledAppRow struct {
	ID           int64     `db:"id"`
	Name         string    `db:"name"`
	CatalogAppID string    `db:"catalog_app_id"`
	CreatedAt    time.Time `db:"created_at"`
}
