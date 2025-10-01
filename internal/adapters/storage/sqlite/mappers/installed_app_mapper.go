package mappers

import (
	"autohost-cli/internal/adapters/storage/sqlite/models"
	"autohost-cli/internal/domain"
)

func toDomainInstalledApp(r models.InstalledAppRow) domain.InstalledApp {
	return domain.InstalledApp{
		ID:           r.ID,
		Name:         r.Name,
		CatalogAppID: r.CatalogAppID,
		CreatedAt:    r.CreatedAt,
	}
}
