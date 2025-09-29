package app

import (
	"autohost-cli/internal/adapters/storage/sqlite"
	"context"
)

type CatalogService struct{ Catalog sqlite.CatalogRepo }

func (s CatalogService) List(ctx context.Context) (any, error) {
	return s.Catalog.ListApps(ctx)
}
