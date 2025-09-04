// internal/services/catalog_service.go
package services

import (
	"autohost-cli/internal/repo"
	"context"
)

type CatalogService struct{ Catalog repo.CatalogRepo }

func (s CatalogService) List(ctx context.Context) (any, error) {
	return s.Catalog.ListApps(ctx)
}
