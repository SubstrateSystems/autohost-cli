package domain

import "time"

type InstalledApp struct {
	ID           int64     `db:"id" `
	Name         string    `db:"name"`
	CatalogAppID string    `db:"catalog_app_id"`
	CreatedAt    time.Time `db:"created_at"`
}

type AppConfig struct {
	Name     string
	Template string
	Port     string
	MySQL    *MySQLConfig
}

type MySQLConfig struct {
	RootPassword string
	User         string
	Password     string
	Database     string
	Port         string
}
