package domain

import (
	"context"
	"time"
)

type CatalogApp struct {
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type InstalledApp struct {
	ID           int64
	Name         string
	CatalogAppID string
	CreatedAt    time.Time
}

// Configuración: pertenece al dominio si expresa reglas del negocio
type AppConfig struct {
	Name     string
	Template string
	Port     string
	MySQL    *MySQLConfig
	Postgres *PostgresConfig
}

type MySQLConfig struct {
	RootPassword string
	User         string
	Password     string
	Database     string
	Port         string
}

type PostgresConfig struct {
	User     string
	Password string
	Database string
	Port     string
}

// ????????????????
type InstalledRepo interface {
	List(ctx context.Context) ([]InstalledApp, error)
	Remove(ctx context.Context, name string) error
	IsInstalledApp(ctx context.Context, name string) (bool, error)
	Add(ctx context.Context, app InstalledApp) error
}

type CatalogItem struct {
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CatalogRepo interface {
	ListApps(ctx context.Context) ([]CatalogItem, error)
}
