package ports

import (
	"autohost-cli/internal/domain"
	"context"
)

type CatalogRepository interface {
	ListApps(ctx context.Context) ([]domain.CatalogApp, error)
	FindByName(ctx context.Context, name domain.AppName) (domain.CatalogApp, error)
}

type InstalledRepository interface {
	List(ctx context.Context) ([]domain.InstalledApp, error)
	Install(ctx context.Context, app domain.InstalledApp) error
	Remove(ctx context.Context, name domain.AppName) error
	IsInstalled(ctx context.Context, name domain.AppName) (bool, error)
}
