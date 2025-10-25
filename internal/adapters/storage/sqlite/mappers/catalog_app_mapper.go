package mappers

import (
	"autohost-cli/internal/adapters/storage/sqlite/models"
	"autohost-cli/internal/domain"
)

func toDomainCatalogApp(r models.CatalogAppRow) domain.CatalogApp {
	return domain.CatalogApp{
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
