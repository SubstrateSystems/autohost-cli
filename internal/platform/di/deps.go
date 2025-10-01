package di

import (
	"autohost-cli/internal/adapters/storage/sqlite"
	"autohost-cli/internal/domain"
	"context"
	"database/sql"
)

type Deps struct {
	DB       *sql.DB
	Repos    Repos
	Services Services
}

type Repos struct {
	Installed *sqlite.InstalledRepo
	Catalog   sqlite.CatalogRepo
}

type Services struct {
	App     AppService
	Catalog CatalogService
}

type AppService struct {
	Installed *sqlite.InstalledRepo
}

func (s AppService) ListInstalled(ctx context.Context) ([]domain.InstalledApp, error) {
	return s.Installed.List(ctx)
}

func (s AppService) RemoveApp(ctx context.Context, name string) error {
	return s.Installed.Remove(ctx, name)
}

func (s AppService) IsAppInstalled(ctx context.Context, name string) (bool, error) {
	return s.Installed.IsInstalledApp(ctx, name)
}

type CatalogService struct{ Catalog sqlite.CatalogRepo }

func (s CatalogService) List(ctx context.Context) (any, error) {
	return s.Catalog.ListApps(ctx)
}
